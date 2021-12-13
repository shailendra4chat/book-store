package db

import (
	"fmt"
	"time"

	"github.com/shailendra4chat/book-store/users/config"
	"github.com/shailendra4chat/book-store/users/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB
var err error

func connectionString() string {
	host := config.Conf("DBHOST")
	user := config.Conf("DBUSER")
	password := config.Conf("DBPASSWORD")
	dbname := config.Conf("UAPP_DBNAME")
	port := config.Conf("DBPORT")
	sslmode := config.Conf("DBSSLMODE")

	//Define DB connection string
	conString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", host, user, password, dbname, port, sslmode)

	return conString
}

func DbConnection() {

	fmt.Println("Connection String==> ", connectionString())
	fmt.Println("Sleeping for 10 seconds...")
	time.Sleep(10 * time.Second)
	fmt.Println("Sleep Over...")

	Database, err = gorm.Open(postgres.Open(connectionString()), &gorm.Config{})

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println("DB Connected!")

	Database.AutoMigrate(&models.User{})
}
