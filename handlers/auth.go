package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dto "chat/DTO/userdto"

	"github.com/dgrijalva/jwt-go"
)

var (
	signKey = os.Getenv("SigningKey")
)

// AuthMiddleWare
func AuthMiddleWare(next func() http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve token from request
		token, err := getToken(r)
		if err != nil {
			// redirecting to login page
			log.Printf("Couldn't get token: %v, err: %v", token, err)
			r.Method = "GET"
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// retrieve user's id from token
		_, err = validateToken(token)
		if err != nil {
			// redirecting to login page
			log.Printf("couldn't validate token, err: %v", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next().ServeHTTP(w, r)
	})
}

// Generating token and send it to user
func AuthUser(w http.ResponseWriter, r *http.Request, user dto.GetUserDTO) error {
	token, err := generateToken(user)
	if err != nil {
		log.Print(err)
		return err
	}
	cookies := &http.Cookie{}
	cookies.Name = "Authorization"
	cookies.Value = "Bearer " + token
	cookies.Path = "/"
	cookies.Expires = time.Now().Add(12 * time.Hour)
	http.SetCookie(w, cookies)
	return nil
}

// generateToken generates jwt token
func generateToken(user dto.GetUserDTO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Authorization": strconv.Itoa(user.Id),
		"exp":           time.Now().Add(12 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		log.Print(err)
		err = fmt.Errorf("server error")
		return "", err
	}

	return tokenStr, nil
}
// validateToken validating token and returning user's id or error
func validateToken(tokenString string) (string, error) {
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
		return "", fmt.Errorf("failed to parse token: %v", err)
	}
	// Extract the userID field from the token's payload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims format")
	}
	userId, ok := claims["Authorization"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid userID field")
	}
	// Convert the userID to an int and return it
	return userId, nil
}

// getToken gets token from cookies
func getToken(r *http.Request) (string, error) {
	token, err := r.Cookie("Authorization")
	if err != nil {
		newErr := fmt.Sprintf("couldn't get cookies, err: %v", err)
		log.Print(newErr)
		return "", err
	}
	splitToken := strings.Split(token.Value, " ")
	if len(splitToken) != 2 {
		err := fmt.Errorf("Token not found, token: %v", token)
		log.Printf("Unable to get token, err: %v", err)
		return "", err
		// Error: Bearer token not in proper format
	}

	return splitToken[1], nil
}
