package middleware

import (
	redisDto "chat/DTO/jwt"
	jwtcontroller "chat/controllers/jwtController"
	"chat/logger"

	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtController = jwtcontroller.New()
	signKey       = os.Getenv("SigningKey")

	jwtRefreshExpiration = 30 * (24 * time.Hour)
	jwtAccessExpiration  = 15 * time.Minute
)

// GenerateToken generates a jwt token with the specified claims and expiration time
func GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		return "", fmt.Errorf("Couldn't sign token: %v", err)
	}

	if id, ok := claims["RefreshToken"].(string); ok {
		err = jwtController.SaveToken(redisDto.RedisDto{Token: tokenStr, Key: "RefreshToken", Id: id, Expiration: jwtRefreshExpiration})
		if err != nil {
			return "", fmt.Errorf("Couldn't save refresh token to Redis: %v", err)
		}
	}
	if id, ok := claims["AccesssToken"].(string); ok {
		err = jwtController.SaveToken(redisDto.RedisDto{Token: tokenStr, Key: "AccessToken", Id: id, Expiration: jwtAccessExpiration})
		if err != nil {
			return "", fmt.Errorf("Couldn't save access token to Redis: %v", err)
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
	return GenerateToken(claims)
}

// GenerateAccessToken generates an access token for the specified user ID
func GenerateAccessToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"AccesssToken": userId,
		"exp":          jwtAccessExpiration,
	}
	return GenerateToken(claims)
}

// DeleteCookies sets cookies by setting new one that will expire immediately
func DeleteCookies(w http.ResponseWriter) {
	accesCookies := &http.Cookie{}
	accesCookies.Name = "AccesssToken"
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
		Name:     "RefreshToken",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(30 * (24 * time.Hour)),
		HttpOnly: true,
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
		Name:     "AccesssToken",
		Value:    accessToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
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

func ValidateToken(cookies *http.Cookie) (jwt.MapClaims, error) {
	tokenString := cookies.Value
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
		accesCookies, err := r.Cookie("AccesssToken")
		if err != nil {
			logger.Error("Couldn't get access accesCookies", err, funcName)
			refreshCookies, err := r.Cookie("RefreshToken")
			if err != nil {
				logger.Error("Coudn't get refresh token from cookies", err, funcName)
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}

			claims, err := ValidateToken(refreshCookies)
			if err != nil {
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}

			userId, ok := claims["RefreshToken"].(string)
			if !ok {
				logger.Error("Couldn't convert claims to int", err, funcName)

				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}

			tokenString := refreshCookies.Value
			redisToken, err := jwtController.GetToken(redisDto.RedisDto{Id: userId, Key: "RefreshToken"})
			if err != nil {
				logger.Error("Couldn't get RefreshToken from redis", err, funcName)
				DeleteCookies(w)
				jwtController.DeleteToken(redisDto.RedisDto{Id: userId, Key: "RefreshToken"})
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}
			if redisToken != tokenString {
				logger.Error("User's refresh token and token from redis don't match", err, funcName)
				DeleteCookies(w)
				jwtController.DeleteToken(redisDto.RedisDto{Id: userId, Key: "RefreshToken"})
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}
			err = SetAccessCookies(w, userId)
			if err != nil {
				logger.Error("Couldn't generate acces token", err, funcName)
				http.Redirect(w, r, "/login/", http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		tokenString := accesCookies.Value
		claims, err := ValidateToken(accesCookies)
		if err != nil {
			logger.Error("Couldn't validate token", err, funcName)

			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		userId, ok := claims["AccesssToken"].(string)
		if !ok {
			logger.Error("Couldn't convert claims to int", err, funcName)

			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		// Getting token from redis by user's id
		redisToken, err := jwtController.GetToken(redisDto.RedisDto{Id: userId, Key: "AccessToken"})
		if err != nil {
			logger.Error("Couldn't get access token from redis", err, funcName)
			DeleteCookies(w)
			jwtController.DeleteToken(redisDto.RedisDto{Id: userId, Key: "AccessToken"})
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		if redisToken != tokenString {
			logger.Error("Token from redis doesn't match", err, funcName)
			DeleteCookies(w)
			jwtController.DeleteToken(redisDto.RedisDto{Id: userId, Key: "AccessToken"})
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("AccesssToken")
	if err != nil {
		_, err := r.Cookie("RefreshToken")
		if err != nil {
			return false
		}
	}
	return true
}

// LogoutHandler deletes user's access and refresh tokens
func LogOut(w http.ResponseWriter, r *http.Request) {
	funcName := logger.GetFuncName()

	// Get the refresh cookies from the request
	cookie, err := r.Cookie("RefreshToken")
	if err != nil {
		logger.Error("Failed to get refresh token cookie", err, funcName)
		cookie, err = r.Cookie("AccesssToken")
		if err != nil {
			logger.Error("Failed to get refresh token cookie", err, funcName)

		}
	}
	userId, err := GetUserIdFromToken(cookie)
	if err != nil {
		logger.Error("Failed to get user id from token", err, funcName)

	}
	// Delete the cookies

	var redis redisDto.RedisDto
	redis.Id = userId
	redis.Key = "RefreshToken"
	// Delete the tokens from Redis
	err = jwtController.DeleteToken(redis)
	if err != nil {
		logger.Error("Failed to delete refresh token from Redis", err, funcName)
	}
	redis.Id = userId
	redis.Key = "AccessToken"
	err = jwtController.DeleteToken(redis)
	if err != nil {
		logger.Error("Failed to delete access token from Redis", err, funcName)
	}
}

// function that gets user id from jwt refresh or access token
func GetUserIdFromToken(cookies *http.Cookie) (string, error) {
	funcName := logger.GetFuncName()
	claims, err := ValidateToken(cookies)
	if err != nil {
		logger.Error("Couldn't validate token", err, funcName)
		return "", err
	}
	userId, ok := claims["RefreshToken"].(string)
	if !ok {
		userId, ok = claims["AccesssToken"].(string)
		if !ok {
			logger.Error("Couldn't convert claims to int", err, funcName)
			return "", err
		}
	}
	return userId, nil
}
