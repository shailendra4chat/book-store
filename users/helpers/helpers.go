package helpers

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// Check If Admin - based on token
func IsAdmin(userToken string) bool {
	token, _ := jwt.Parse(userToken, nil)
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["Admin"] == true
}

// Get current user name
func CurrentUser(userToken string) string {
	token, _ := jwt.Parse(userToken, nil)
	claims, _ := token.Claims.(jwt.MapClaims)

	name := fmt.Sprint(claims["Name"])

	return name
}
