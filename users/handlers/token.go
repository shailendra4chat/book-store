package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shailendra4chat/book-store/users/helpers"
)

// ValidateToken ... Validate user token
// @Summary Validate user token
// @Description Validate user token
// @Tags Token
// @Success 200 {array} models.User
// @Failure 401 {object} object
// @Router /auth/token [get]
// @Security ApiKeyAuth
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(true)
}

// ValidateAdminToken ... Validate admin token
// @Summary Validate admin token
// @Description Validate admin token
// @Tags Token
// @Success 200 {array} models.User
// @Failure 401 {object} object
// @Router /auth/token/admin [get]
// @Security ApiKeyAuth
func ValidateAdminToken(w http.ResponseWriter, r *http.Request) {

	isAdmin := helpers.IsAdmin(r.Header.Get("x-access-token"))

	if isAdmin {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Admin Access")
	} else {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Admin access not found!")
	}
}
