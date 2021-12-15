package main

import (
	"github.com/shailendra4chat/book-store/books/db"
	"github.com/shailendra4chat/book-store/books/routes"
)

// @title Books API documentation
// @version 1.0.0
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-access-token

func main() {
	// DB Connection
	db.DbConnection()

	// Initialise routing
	routes.HandleRouting()
}
