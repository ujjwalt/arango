package arango

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Connect("http://localhost:8529")
	os.Exit(m.Run())
}

func TestCurrentDB(t *testing.T) {
	t.Log("Test if we can fetch the current database")
	db, err := CurrentDB()
	if err != nil {
		t.Fatal(err)
	}

	if len(db.Name) < 1 {
		t.Error("db does not have a name")
	}

	if len(db.Id) < 1 {
		t.Error("db does not have a id")
	}

	if len(db.Path) < 1 {
		t.Error("db does not have a path")
	}
	log.Println(db)
}
