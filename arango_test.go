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

func TestCreateDocument(t *testing.T) {
	t.Log("Test if we can create a new document")
	doc := &map[string]interface{}{
		"name":  "Ujjwal Thaakar",
		"email": "ujjwal@gmail.com",
	}
	d, err := CreateDocument(*doc, "ujjwal", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(d)

	if _, ok := d["_id"].(string); !ok {
		t.Fatal("Document _id should be a string")
	}

	if _, ok := d["_key"].(string); !ok {
		t.Fatal("Document _key should be a string")
	}

	if _, ok := d["_rev"].(string); !ok {
		t.Fatal("Document _rev should be a string")
	}

	log.Println(d)
}

func TestFind(t *testing.T) {
	t.Log("Test if we can find a document by it's handle")
	doc := &map[string]interface{}{
		"name":  "Ujjwal Thaakar",
		"email": "ujjwal@gmail.com",
	}
	d, err := CreateDocument(*doc, "ujjwal", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)

	handle := d["_id"].(string)
	same, err := Find(handle)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(same)

	if d["_id"] != same["_id"] {
		t.Fatalf("Documnt with handle %s not fetched: %s != %s", handle, same["_id"], handle)
	}

	log.Printf("%s: %v", handle, same)
}
