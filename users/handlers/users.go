package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/shailendra4chat/book-store/users/db"
	"github.com/shailendra4chat/book-store/users/helpers"
	"github.com/shailendra4chat/book-store/users/models"
)

// Handle user registration
func Register(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	alreadyRegistered := checkIfRegistered(user.Email)

	if alreadyRegistered {
		http.Error(w, "User already registered with this Email!", http.StatusConflict)
	} else {
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		user.Password = string(pass)

		db.Database.Create(user)

		var resp = map[string]interface{}{"message": "User created successfully!"}

		json.NewEncoder(w).Encode(resp)
	}
}

func checkIfRegistered(email string) bool {
	user := &models.User{}

	if err := db.Database.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = false
		return resp
	}
	return true
}

// Handle user login
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	resp := checkCreds(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

func checkCreds(email, password string) map[string]interface{} {
	user := &models.User{}

	if err := db.Database.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = map[string]interface{}{"message": "Email address not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 100).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &models.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Admin:  user.Admin,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"message": "logged in"}
	user.Password = "*****"
	resp["token"] = tokenString //Store the token in the response
	resp["user"] = user

	delete(resp, "password")
	return resp
}

// Handle all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.Database.Find(&users).Error; err != nil {
		fmt.Println(err)
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// Handle user update
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	if err := db.Database.Where("id = ?", id).First(user).Error; err != nil {
		fmt.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(user)
	db.Database.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

// Handle user delete - Only admin can delete
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	isAdmin := helpers.IsAdmin(r.Header.Get("x-access-token"))
	if isAdmin {
		user := &models.User{}
		params := mux.Vars(r)
		var id = params["id"]
		if err := db.Database.Where("id = ?", id).First(user).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewDecoder(r.Body).Decode(user)
		db.Database.Delete(&user)
		json.NewEncoder(w).Encode("User deleted")
	} else {
		http.Error(w, "Not authorised", http.StatusUnauthorized)
		return
	}
}
