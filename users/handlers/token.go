package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shailendra4chat/book-store/users/helpers"
)

// Validate user token
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(true)
}

// Validate admin token
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
