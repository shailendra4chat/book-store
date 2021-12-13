package main

import (
	"github.com/shailendra4chat/book-store/books/db"
	"github.com/shailendra4chat/book-store/books/routes"
)

func main() {
	// DB Connection
	db.DbConnection()

	// Initialise routing
	routes.HandleRouting()
}
