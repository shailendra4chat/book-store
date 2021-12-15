package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shailendra4chat/book-store/books/db"
	"github.com/shailendra4chat/book-store/books/helpers"
	"github.com/shailendra4chat/book-store/books/models"
)

// AddBook ... Add Book
// @Summary Add Book
// @Description Add Book
// @Tags Books
// @Accept json
// @Param book body models.Book true "Book Data"
// @Success 200 {object} models.Book
// @Failure 400,500 {object} object
// @Router /book [post]
// @Security ApiKeyAuth
func AddBook(w http.ResponseWriter, r *http.Request) {

	isAdmin := helpers.IsAdmin(r.Header.Get("x-access-token"))

	if isAdmin {
		books := &models.Book{}
		json.NewDecoder(r.Body).Decode(books)

		alreadyRegistered := checkIfAdded(books.Title)

		if alreadyRegistered {
			http.Error(w, "Book already added with this Title!", http.StatusConflict)
		} else {
			db.Database.Create(books)
			var resp = map[string]interface{}{"message": "Book added successfully!"}

			json.NewEncoder(w).Encode(resp)
		}

	} else {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
	}
}

func checkIfAdded(title string) bool {
	books := &models.Book{}

	if err := db.Database.Where("Title = ?", title).First(books).Error; err != nil {
		var resp = false
		return resp
	}
	return true
}

// GetBooks ... Get all books
// @Summary Get all books
// @Description get all books
// @Tags Books
// @Success 200 {array} models.Book
// @Failure 404 {object} object
// @Router /books [get]
// @Security ApiKeyAuth
func GetBooks(w http.ResponseWriter, r *http.Request) {

	isAuthorisedUser := helpers.IsAuthorisedUser(r.Header.Get("x-access-token"))

	if isAuthorisedUser {
		var books []models.Book
		if err := db.Database.Find(&books).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "Books not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(books)

	} else {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
	}
}

// UpdateBook ... Update Book by Id
// @Summary Update book by Id
// @Description Update book by Id
// @Tags Books
// @Accept json
// @Param id path string true "Book ID"
// @Param book body models.Book true "Book Data"
// @Success 201 {object} models.Book
// @Failure 400,500 {object} object
// @Router /book/{id} [put]
// @Security ApiKeyAuth
func UpdateBook(w http.ResponseWriter, r *http.Request) {

	isAdmin := helpers.IsAdmin(r.Header.Get("x-access-token"))

	if isAdmin {
		book := &models.Book{}
		params := mux.Vars(r)
		var id = params["id"]
		if err := db.Database.Where("id = ?", id).First(book).Error; err != nil {
			fmt.Println(err)
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		json.NewDecoder(r.Body).Decode(book)
		db.Database.Save(&book)
		json.NewEncoder(w).Encode(&book)

	} else {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
	}
}
