package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shailendra4chat/book-store/books/config"
	"github.com/shailendra4chat/book-store/books/handlers"

	_ "github.com/shailendra4chat/book-store/books/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func HandleRouting() {
	host := ":" + config.Conf("BAPP_PORT")
	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/book", handlers.AddBook).Methods("POST")

	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")

	r.HandleFunc("/book/{id}", handlers.UpdateBook).Methods("PUT")

	fmt.Printf("Books app is running on port: %v", config.Conf("BAPP_PORT"))
	log.Fatal(http.ListenAndServe(host, r))
}
