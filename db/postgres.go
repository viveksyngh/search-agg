package db

import (
	"database/sql"
	"fmt"
	"os"

	//Postgres database driver
	_ "github.com/lib/pq"
)

//GetEnv get environment variable with default value
func GetEnv(name string, value string) string {
	envValue := os.Getenv(name)
	if len(envValue) == 0 {
		return value
	}
	return envValue
}

//Connection creates and returns database connection
func Connection() (*sql.DB, error) {
	dbUser := GetEnv("DB_USER", "viveks")
	dbPassword := GetEnv("DB_PASSWORD", "")
	dbPort := GetEnv("DB_PORT", "5432")
	dbName := GetEnv("DB_NAME", "searchdb")
	dbHost := GetEnv("DB_HOST", "localhost")

	var connectionString = "sslmode=disable host=%s port=%s user=%s " +
		"dbname=%s password=%s"

	connectionString = fmt.Sprintf(connectionString, dbHost, dbPort, dbUser, dbName, dbPassword)
	fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err.Error())
	}

	return db, err
}
