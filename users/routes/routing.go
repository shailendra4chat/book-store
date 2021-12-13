package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shailendra4chat/book-store/users/config"
	"github.com/shailendra4chat/book-store/users/handlers"
	"github.com/shailendra4chat/book-store/users/middlewares"
)

func HandleRouting() {
	host := ":" + config.Conf("UAPP_PORT")
	r := mux.NewRouter()

	r.Use(middlewares.HeadersMiddleware)

	r.HandleFunc("/register", handlers.Register).Methods("POST")

	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middlewares.JwtVerify)

	s.HandleFunc("/users", handlers.GetUsers).Methods("GET")

	s.HandleFunc("/update-user/{id}", handlers.UpdateUser).Methods("PUT")

	s.HandleFunc("/delete-user/{id}", handlers.DeleteUser).Methods("DELETE")

	// Validate token
	s.HandleFunc("/token", handlers.ValidateToken).Methods("GET")

	// Validate token for Admin
	s.HandleFunc("/token/admin", handlers.ValidateAdminToken).Methods("GET")

	fmt.Printf("Users app is running on port: %v", config.Conf("PORT"))

	log.Fatal(http.ListenAndServe(host, r))
}
