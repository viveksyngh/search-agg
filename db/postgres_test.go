package db

import (
	"fmt"
	"testing"
)

func Test_Connection(t *testing.T) {
	db, err := Connection()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Connected")
}
