package main

import (
	"github.com/shailendra4chat/book-store/users/db"
	"github.com/shailendra4chat/book-store/users/routes"
)

// @title Users API documentation
// @version 1.0.0
// @host localhost:8080
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
