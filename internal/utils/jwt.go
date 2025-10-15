package utils

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken creates a JWT for a given user ID.
func GenerateToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // expires in 1 day
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken validates a JWT and returns the user ID.
func ParseToken(tokenString string) (int, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecret, nil
    })

    if err != nil || !token.Valid {
        return 0, errors.New("invalid token")
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if uid, ok := claims["user_id"].(float64); ok {
            return int(uid), nil
        }
    }

    return 0, errors.New("invalid token claims")
}