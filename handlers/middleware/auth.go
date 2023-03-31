package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	redisDto "chat/DTO/redis_jwt"
	rediscontroller "chat/controllers/redisController"
	"chat/entities/redis_jwt"
	"chat/logger"

	"github.com/dgrijalva/jwt-go"
)

var (
	redisDb         = redis_jwt.RedisData{}
	redisController = rediscontroller.New(&redisDb)
	signKey         = os.Getenv("SigningKey")

	jwtRefreshExpiration = time.Now().Add(30 * (24 * time.Hour))
	jwtAccessExpiration  = time.Now().Add(15 * time.Minute)
)

type Token struct {
	Type     string
	StrToken string
	Value    *jwt.Token
}

// GenerateToken generates a jwt token with the specified claims and expiration time
func GenerateToken(claims jwt.MapClaims, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		return "", fmt.Errorf("Couldn't sign token: %v", err)
	}
	if id, ok := claims["Authorization"].(string); ok {
		err = redisController.SaveToken(redisDto.RedisDto{Token: tokenStr, Key: "jwt", Id: id})
		if err != nil {
			return "", fmt.Errorf("Couldn't save token to Redis: %v", err)
		}
	} else if id, ok := claims["RefreshToken"].(string); ok {
		err = redisController.SaveToken(redisDto.RedisDto{Token: tokenStr, Key: "RefreshToken", Id: id})
		if err != nil {
			return "", fmt.Errorf("Couldn't save token to Redis: %v", err)
		}
	}
	return tokenStr, nil
}

// GenerateRefreshToken generates a refresh token for the specified user ID
func GenerateRefreshToken(userId string) (string, error) {

	claims := jwt.MapClaims{
		"RefreshToken": userId,
		"exp":          jwtRefreshExpiration,
	}
	return GenerateToken(claims, jwtRefreshExpiration)
}

// GenerateAccessToken generates an access token for the specified user ID
func GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"Authorization": userId,
		"exp":           jwtAccessExpiration,
	}

	return GenerateToken(claims, jwtAccessExpiration)
}

// DeleteCookies sets cookies by setting new one that will expire immediately
func DeleteCookies(w http.ResponseWriter) {
	accesCookies := &http.Cookie{}
	accesCookies.Name = "Authorization"
	accesCookies.Expires = time.Now().Add(-1 * time.Hour)
	accesCookies.Path = "/"
	accesCookies.Domain = "localhost"

	refreshCookies := &http.Cookie{}
	refreshCookies.Name = "RefreshToken"
	refreshCookies.Expires = time.Now().Add(-1 * time.Hour)
	refreshCookies.Path = "/"
	refreshCookies.Domain = "localhost"
	http.SetCookie(w, accesCookies)
	http.SetCookie(w, refreshCookies)
}

// SetRefreshCookies sets refresh token to user's cookies
func SetRefreshCookies(w http.ResponseWriter, userId string) error {
	refreshToken, err := GenerateRefreshToken(userId)
	if err != nil {
		return err
	}
	refreshCookie := &http.Cookie{
		Name:    "RefreshToken",
		Value:   "Bearer " + refreshToken,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(30 * (24 * time.Hour)),
	}
	http.SetCookie(w, refreshCookie)
	return nil
}

func SetAccessCookies(w http.ResponseWriter, userId string) error {
	accessToken, err := GenerateAccessToken(userId)
	if err != nil {
		return err
	}
	accessCookie := &http.Cookie{
		Name:    "Authorization",
		Value:   "Bearer " + accessToken,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(15 * time.Minute),
	}
	http.SetCookie(w, accessCookie)
	return nil
}

// AuthUser generates token and send it to user's cookies
func AuthUser(w http.ResponseWriter, userId string) error {
	funcName := logger.GetFuncName()
	err := SetAccessCookies(w, userId)
	if err != nil {
		logger.Error("Couldn't set acces token to cookies", err, funcName)
		return err
	}
	err = SetRefreshCookies(w, userId)
	if err != nil {
		logger.Error("Couldn't set refresh token to cookies", err, funcName)
		return err
	}
	return nil
}

func validateToken(cookies *http.Cookie) (jwt.MapClaims, error) {
	tokenString := strings.Split(cookies.Value, " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return jwt.MapClaims{}, err
	}
	return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		funcName := logger.GetFuncName()
		accesCookies, err := r.Cookie("Authorization")
		if err != nil {
			logger.Error("Couldn't get access accesCookies", err, funcName)
			refreshCookies, err := r.Cookie("RefreshToken")

			if err != nil {
				logger.Error("Coudn't get refresh token from cookies", err, funcName)
				http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}
			claims, err := validateToken(refreshCookies)
			if err != nil {

				http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}
			userId, ok := claims["RefreshToken"].(string)
			if !ok {
				logger.Error("Couldn't convert claims to int", err, funcName)

				http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}

			err = SetAccessCookies(w, userId)
			if err != nil {
				logger.Error("Couldn't generate acces token", err, funcName)

				http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// tokenString := strings.Split(accesCookies.Value, " ")[1]
		claims, err := validateToken(accesCookies)
		if err != nil {
			logger.Error("Couldn't validate token", err, funcName)

			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		userId, ok := claims["Authorization"].(string)
		if !ok {
			logger.Error("Couldn't convert claims to int", err, funcName)

			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		// redistoken, err := redisController.GetJWT(strconv.Itoa(userId))
		// if err != nil {
		// 	logger.Error("Couldn't get jwt from redis", err, funcName)
		//
		// 	http.Redirect(w, r, "/login/", http.StatusFound)
		// 	return
		// }
		// if redistoken != tokenString {
		// 	logger.Error("Token from redis doesn't match", err, funcName)
		//
		// 	http.Redirect(w, r, "/login/", http.StatusFound)
		// 	return
		// }
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("Authorization")
	if err != nil {
		_, err := r.Cookie("RefreshToken")
		if err != nil {
			return false
		}
	}
	return true
}
