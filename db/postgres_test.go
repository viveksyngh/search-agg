package db

import (
	"testing"
)

func Test_Connection(t *testing.T) {
	_, err := Connection()
	if err != nil {
		t.Fatal(err.Error())
	}
}
