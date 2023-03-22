package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	rdto "chat/DTO/redis_jwt"
	rediscontroller "chat/controllers/redisController"
	"chat/entities/redis_jwt"
	"chat/logger"

	"github.com/dgrijalva/jwt-go"
)

var (
	rdbService    = redis_jwt.UserJwt{}
	rdbController = rediscontroller.New(&rdbService)
	signKey       = os.Getenv("SigningKey")
)

// AuthMiddleWare checks if user is authorized if not - redirects to login page
func AuthMiddleWare(next func() http.Handler) http.Handler {
	funcName := logger.GetFuncName()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve token from request
		token, err := GetToken(r)
		if err != nil || token == "" {
			// redirecting to login page
			logger.Error("Couldn't get token", err, funcName)
			DeleteCookies(w)
			r.Method = "GET"
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next().ServeHTTP(w, r)
	})
}

// DeleteCookies sets cookies by setting new one that will expire immediately
func DeleteCookies(w http.ResponseWriter) {
	cookies := &http.Cookie{}
	cookies.Name = "Authorization"
	cookies.Expires = time.Unix(0, 0)
	cookies.Path = "/"
	cookies.Domain = "localhost"
	http.SetCookie(w, cookies)
}

// AuthUser generates token and send it to user
func AuthUser(w http.ResponseWriter, userId int) error {
	funcName := logger.GetFuncName()
	token, err := GenerateToken(userId)
	if err != nil {
		logger.Error("Couldn't generate token", err, funcName)
		return err
	}
	cookies := &http.Cookie{}
	cookies.Name = "Authorization"
	cookies.Value = "Bearer " + token
	cookies.Path = "/"
	cookies.Domain = "localhost"
	cookies.Expires = time.Now().Add(15 * time.Minute)
	http.SetCookie(w, cookies)
	return nil
}

// generateRefreshToken generates refresh token for user
func generateRefreshToken(userId int) (string, error) {
	funcName := logger.GetFuncName()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": strconv.Itoa(userId),
		"exp":           time.Now().Add(30 * (24 * time.Hour)).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		logger.Error("Couldn't sign refresh token", err, funcName)
		return "", err
	}
	return tokenStr, nil

}

// GenerateToken generates jwt token
func GenerateToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": strconv.Itoa(userId),
		"exp":           time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		funcName := logger.GetFuncName()
		logger.Error("Couldn't sign token", err, funcName)
		err = fmt.Errorf("server error")
		return "", err
	}
	var redisData rdto.RdbDTO
	redisData.Token = tokenStr
	redisData.Id = strconv.Itoa(userId)
	rdbController.SaveJWT(redisData)
	return tokenStr, nil
}

// ValidateToken validating token and returning user's id or error
func ValidateToken(tokenString string) (int, error) {
	funcName := logger.GetFuncName()
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
		return -1, fmt.Errorf("failed to parse token: %v", err)
	}
	userId, err := GetUserIdFromToken(*token)
	if err != nil {
		logger.Error("couldn't get user's id from token", err, funcName)
		return -1, err
	}
	userData := rdto.RdbDTO{}
	err, equal := rdbController.CompareJWT(userData)
	if err != nil {
		return -1, fmt.Errorf("token is not valid, err: %v", err)
	}
	if !equal {
		return -1, fmt.Errorf("token is not equal")
	}
	// Convert the userID to an int and return it
	intUserId, _ := strconv.Atoi(userId)
	return intUserId, nil
}

func GetUserIdFromToken(token jwt.Token) (string, error) {
	// Extract the userID field from the token's payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims format")
	}
	userId, ok := claims["Authorization"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid userID field")
	}
	return userId, nil
}

// GetToken gets token from cookies
func GetToken(r *http.Request) (string, error) {
	token, err := r.Cookie("Authorization")
	if err != nil {
		return "", err
	}
	splitToken := strings.Split(token.Value, " ")
	if len(splitToken) != 2 || splitToken[1] == "" {
		err := fmt.Errorf("token not found, token: %v", token)
		return "", err
	}
	return splitToken[1], nil
}

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(w http.ResponseWriter, r *http.Request) {
	token, err := GetToken(r)
	if err != nil || token == "" {
		DeleteCookies(w)
		return
	}
	http.Redirect(w, r, "/chat/", http.StatusFound)
}
