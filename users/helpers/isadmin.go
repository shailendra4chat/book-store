package helpers

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// Check If Admin - based on token
func IsAdmin(userToken string) bool {
	token, _ := jwt.Parse(userToken, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["Admin"] == true
}
