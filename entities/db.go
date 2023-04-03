package entity

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	l   = logrus.New()
	ctx = context.Background()
)

// DbConnect establishes a connection to the database.
func MySQLConnect() (*sql.DB, error) {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dataSourceName := os.Getenv("DB_SOURCE")
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return db, err
	}
	return db, nil
}

// DbTablesInit creates massage and users tables in the database if they do not already exist.
func MySQLTablesInit() error {
	db, err := MySQLConnect()
	if err != nil {
		return err
	}
	// Creating tables
	_, err = db.Query("CREATE TABLE IF NOT EXISTS chat.messages (" +
		"id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, " +
		"user_name VARCHAR(16), " +
		"user_id INT," +
		"text VARCHAR(4096), " +
		"created_at DATETIME, " +
		"updated_at DATETIME" +
		");")
	if err != nil {
		log.Printf("Couldn't create table 'messages'")
		return err
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS chat.users (" +
		"id INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
		"name VARCHAR(20)," +
		"email VARCHAR(20)," +
		"password VARCHAR(200)" +
		");")
	if err != nil {
		log.Printf("Couldn't create table 'users'")
		return err
	}
	return nil
}

func RedisConnect() (*redis.Client, error) {
	dbNum, err := strconv.Atoi(os.Getenv("Redis_DB"))
	if err != nil {
		log.Fatalf("Error converting RedisDB to integer: %v", err)
	}
	redisInstance := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddr"),
		Username: os.Getenv("RedisUsername"),
		Password: os.Getenv("RedisPassword"),
		DB:       dbNum,
	})
	status := redisInstance.Ping(ctx)
	if status.Err() != nil {
		return redisInstance, fmt.Errorf("Couldn't connect to redis, err:%v", status.Err())
	}
	return redisInstance, nil
}
