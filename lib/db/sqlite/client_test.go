package sqlite3

import "testing"

func TestConnect(t *testing.T) {
	DSN = "./test.db"

	err := Connect()
	if err != nil {
		t.Fatal(err)
	}
}
