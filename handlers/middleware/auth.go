package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dto "chat/DTO/userdto"
	logger "chat/logger"

	"github.com/dgrijalva/jwt-go"
)

var (
	signKey = os.Getenv("SigningKey")
)

// AuthMiddleWare
func AuthMiddleWare(next func() http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve token from request
		token, err := GetToken(r)
		if err != nil || token == "" {
			// redirecting to login page
			funcName := logger.GetFuncName()
			logger.Error("Couldn't get token", err, funcName)
			r.Method = "GET"
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next().ServeHTTP(w, r)
	})
}

// Generating token and send it to user
func AuthUser(w http.ResponseWriter, r *http.Request, user dto.GetUserDTO) error {
	token, err := GenerateToken(user)
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't generate token", err, funcName)
		return err
	}
	cookies := &http.Cookie{}
	cookies.Name = "Authorization"
	cookies.Value = "Bearer " + token
	cookies.Path = "/"
	cookies.Expires = time.Now().Add(15 * time.Minute)
	http.SetCookie(w, cookies)
	return nil
}

// GenerateToken generates jwt token
func GenerateToken(user dto.GetUserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": strconv.Itoa(user.Id),
		"exp":           time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Coudn't sign token", err, funcName)
		err = fmt.Errorf("server error")
		return "", err
	}

	return tokenStr, nil
}

// validateToken validating token and returning user's id or error
func ValidateToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad signing method: %v", token.Header["alg"])
		}
		// Return the secret key used to sign the token
		return []byte(signKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}
	// Extract the userID field from the token's payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("unexpected claims format")
	}
	userId, ok := claims["Authorization"].(string)
	if !ok {
		return 0, fmt.Errorf("missing or invalid userID field")
	}
	// Convert the userID to an int and return it
	intUserId, _ := strconv.Atoi(userId)
	return intUserId, nil
}

// GetToken gets token from cookies
func GetToken(r *http.Request) (string, error) {
	token, err := r.Cookie("Authorization")
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("couldn't get cookies", err, funcName)
		return "", err
	}
	splitToken := strings.Split(token.Value, " ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		err := fmt.Errorf("Token not found, token: %v", token)
		funcName := logger.GetFuncName()
		logger.Error("Unable to get token", err, funcName)
		return "", err
	}

	return splitToken[1], nil
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	token, err := GetToken(r)
	if err != nil || token == "" {
		return
	}
	http.Redirect(w, r, "/chat/", http.StatusFound)
}
