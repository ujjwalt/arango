package arango

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	NewConnection("http://localhost:8529")
	os.Exit(m.Run())
}

func TestCurrentDatabase(t *testing.T) {
	db, err := CurrentDatabase()
	if err != nil || db == nil {
		t.Error(err)
	}
	fmt.Println(*db)
}

func TestAccessibleDatabases(t *testing.T) {
	dbs, err := AccessibleDatabases()
	if err != nil || len(dbs) < 1 {
		t.Error(err)
	}
	for _, db := range dbs {
		fmt.Println(db)
	}
}

func TestDatabases(t *testing.T) {
	dbs, err := Databases()
	if err != nil || len(dbs) < 1 {
		t.Error(err)
	}
	for _, db := range dbs {
		fmt.Println(db)
	}
}
