package entity

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DbConnect establishes a connection to the database.
func DbConnect() (*sql.DB, error) {
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

// DbTableInit creates massage and users tables in the database if they do not already exist.
func DbTablesInit() error {
	db, err := DbConnect()
	if err != nil {
		log.Printf("couldn't connect to db, err: %v", err)
		return err
	}
	// Creating tables
	_, err = db.Query("CREATE TABLE IF NOT EXISTS message (" +
		"id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, " +
		"user_name VARCHAR(16), " +
		"text VARCHAR(4096), " +
		"created_at DATETIME, " +
		"updated_at DATETIME" +
		");")
	if err != nil {
		log.Printf("Couldn't create table 'message'")
		return err
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS users (" +
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
