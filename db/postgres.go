package db

import (
	"database/sql"
	"fmt"
	"os"

	//Postgres database driver
	_ "github.com/lib/pq"
)

//Connection creates and returns database connection
func Connection() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	var connectionString = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	connectionString = fmt.Sprintf(connectionString, dbHost, dbPort, dbUser, dbPassword,
		dbName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err.Error())
	}

	return db, err
}
