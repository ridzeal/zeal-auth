package auth

import (
	"fmt"
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key")

func GenerateToken(username string) (string, error) {
	// Check for invalid characters first, before trimming
    if strings.Contains(username, "\n") || strings.Contains(username, "\r") {
        return "", fmt.Errorf("username contains invalid characters")
    }

    // Check for potentially malicious content
    if strings.Contains(username, "<script>") || strings.Contains(username, "';") {
        return "", fmt.Errorf("username contains invalid characters")
    }

    // Trim spaces and validate emptiness after other validations
    username = strings.TrimSpace(username)
    if username == "" {
        return "", fmt.Errorf("username cannot be empty")
    }

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = username
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}