package testdata

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	txdb "github.com/DATA-DOG/go-txdb"

	"github.com/alextanhongpin/go-microservice/infrastructure/database"
)

func init() {
	txdb.Register("mysqltx", "mysql", database.Option{
		User: os.Getenv("TEST_DB_USER"),
		Pass: os.Getenv("TEST_DB_PASS"),
		Name: os.Getenv("TEST_DB_NAME"),
		Host: os.Getenv("TEST_DB_HOST"),
	}.ConnectionString())
}

func NewDB() *sql.DB {
	connName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
	db, err := sql.Open("mysqltx", connName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
