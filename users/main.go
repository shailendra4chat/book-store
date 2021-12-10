package main

import (
	"github.com/shailendra4chat/book-store/users/db"
	"github.com/shailendra4chat/book-store/users/routes"
)

func main() {
	// DB Connection
	db.DbConnection()

	// Initialise routing
	routes.HandleRouting()
}
