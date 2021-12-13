package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shailendra4chat/book-store/books/config"
	"github.com/shailendra4chat/book-store/books/handlers"
	"github.com/shailendra4chat/book-store/books/middlewares"
)

func HandleRouting() {
	host := ":" + config.Conf("PORT")
	r := mux.NewRouter()

	r.Use(middlewares.HeadersMiddleware)

	r.HandleFunc("/book", handlers.AddBook).Methods("POST")

	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")

	r.HandleFunc("/book/{id}", handlers.UpdateBook).Methods("PUT")

	fmt.Printf("Books app is running on port: %v", config.Conf("PORT"))
	log.Fatal(http.ListenAndServe(host, r))
}
