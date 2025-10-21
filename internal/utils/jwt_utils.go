package utils

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID uint, role string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func ParseToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return int64(0), "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		switch expVal := claims["exp"].(type) {
		case float64:
			if float64(time.Now().Unix()) > expVal {
				return int64(0), "", jwt.ErrTokenExpired
			}
		case json.Number:
			expInt, err := expVal.Int64()
			if err != nil {
				return int64(0), "", jwt.ErrTokenInvalidClaims
			}
			if time.Now().Unix() > expInt {
				return int64(0), "", jwt.ErrTokenExpired
			}
		case string:
			expInt, err := strconv.ParseInt(expVal, 10, 64)
			if err != nil {
				return int64(0), "", jwt.ErrTokenInvalidClaims
			}
			if time.Now().Unix() > expInt {
				return int64(0), "", jwt.ErrTokenExpired
			}
		default:
			return int64(0), "", jwt.ErrTokenInvalidClaims
		}

		// Extract user ID
		var userID int64
		if subVal, ok := claims["sub"]; ok {
			switch v := subVal.(type) {
			case float64:
				userID = int64(v)
			case int64:
				userID = v
			case int:
				userID = int64(v)
			case json.Number:
				n, err := v.Int64()
				if err != nil {
					return int64(0), "", jwt.ErrTokenInvalidClaims
				}
				userID = n
			case string:
				n, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return int64(0), "", jwt.ErrTokenInvalidClaims
				}
				userID = n
			default:
				return int64(0), "", jwt.ErrTokenInvalidClaims
			}
		} else {
			return int64(0), "", jwt.ErrTokenInvalidClaims
		}

		// Extract role
		var role string
		if roleVal, ok := claims["role"]; ok {
			switch v := roleVal.(type) {
			case string:
				role = v
			default:
				return int64(0), "", jwt.ErrTokenInvalidClaims
			}
		} else {
			return int64(0), "", jwt.ErrTokenInvalidClaims
		}

		return userID, role, nil
	}

	return int64(0), "", jwt.ErrTokenInvalidClaims
}
