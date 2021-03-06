package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/shailendra4chat/book-store/users/db"
	"github.com/shailendra4chat/book-store/users/helpers"
	"github.com/shailendra4chat/book-store/users/models"
)

// Register ... Register User
// @Summary Register new user based on paramters
// @Description Register new user
// @Tags Users
// @Accept json
// @Param user body models.User true "User Data"
// @Success 200 {object} object
// @Failure 400,500 {object} object
// @Router /register [post]
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

// Login ... Login User
// @Summary Login user based on paramters
// @Description Login user
// @Tags Users
// @Accept json
// @Param user body models.User true "User Data"
// @Success 200 {object} object
// @Failure 400,500 {object} object
// @Router /login [post]
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

// GetUsers ... Get all users
// @Summary Get all users
// @Description get all users
// @Tags Users
// @Success 200 {array} models.User
// @Failure 404 {object} object
// @Router /auth/users [get]
// @Security ApiKeyAuth
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.Database.Find(&users).Error; err != nil {
		fmt.Println(err)
		http.Error(w, "Users not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// UpdateUser ... Update User by Id
// @Summary Update user by Id
// @Description Update user by Id
// @Tags Users
// @Accept json
// @Param id path string true "User ID"
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 400,500 {object} object
// @Router /auth/update-user/{id} [put]
// @Security ApiKeyAuth
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

// DeleteUser ... Delete User by Id
// @Summary Delete user by Id
// @Description Delete user by Id
// @Tags Users
// @Accept json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400,500 {object} object
// @Router /auth/delete-user/{id} [delete]
// @Security ApiKeyAuth
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

// Upload ... Upload file
// @Summary Upload file
// @Description Upload file
// @Tags Users
// @Accept  multipart/form-data
// @Produce  json
// @Param myFile formData file true  "Upload file"
// @Success 200 {string} string "ok"
// @Failure 400,500 {object} object
// @Router /auth/upload [post]
// @Security ApiKeyAuth
func Upload(w http.ResponseWriter, r *http.Request) {

	// Parse our multipart form
	// 10 << 20 specifies a maximum upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Get path
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	currentUser := helpers.CurrentUser(r.Header.Get("x-access-token"))

	tempFile, err := ioutil.TempFile(filepath.Join(path, "thumbs"), currentUser+"*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
