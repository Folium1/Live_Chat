package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
			return
		}
		// retrieve user's id from token
		userId, err := validateToken(token)
		if err != nil {
			// redirecting to login page
			log.Printf("couldn't validate token, err: %v", err)
			http.Redirect(w, r, "/login/", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userId)

		// Call the next handler with the modified request context
		next().ServeHTTP(w, r.WithContext(ctx))
	})
}

// Generating token and send it to user
func AuthUser(w http.ResponseWriter, user dto.GetUserDTO) error {
	token, err := generateToken(user)
	if err != nil {
		log.Print(err)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Value:   token,
		Expires: time.Now().Add(12 * time.Hour),
	})
	return nil
}

func generateToken(user dto.GetUserDTO) (string, error) {
	fmt.Println(user.Id)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(user.Id),
		"exp":    time.Now().Add(12 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		log.Print(err)
		err = fmt.Errorf("server error")
		return "", err
	}
	return tokenStr, nil
}

func validateToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
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
	userID, ok := claims["userId"].(int)
	if !ok {
		return 0, fmt.Errorf("missing or invalid userID field")
	}
	// Convert the userID to an int and return it
	fmt.Printf("\nUser's id: %v\n", userID)
	return userID, nil
}

// Get token from cookies
func getToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Couldn't parse token, err:", err)
		return "", err
	}
	return cookie.Value, nil
}
