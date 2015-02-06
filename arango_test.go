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

func createDoc(t *testing.T) Document {
	doc := &map[string]interface{}{
		"name":  "Ujjwal Thaakar",
		"email": "ujjwal@gmail.com",
	}
	d, err := CreateDocument(*doc, "ujjwal", true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(d)
	return d
}

func TestCreateDocument(t *testing.T) {
	t.Log("Test if we can create a new document")

	d := createDoc(t)
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

	d := createDoc(t)
	handle := d.Id()
	same, err := Find(handle)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(same)

	if d.Rev() != same.Rev() {
		t.Fatalf("Documnt with handle %s not fetched: %s != %s", handle, d.Id(), same.Rev())
	}

	log.Printf("%s: %v", handle, same)
}

func TestFindIf(t *testing.T) {
	t.Log("Test if we can find a document by it's handle and by conditionally matching it's etag")
	doc := &map[string]interface{}{
		"name":  "Ujjwal Thaakar",
		"email": "ujjwal@gmail.com",
	}
	d, err := CreateDocument(*doc, "ujjwal", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)

	handle := d.Id()
	etag := d.Rev()

	same, err := FindIf(handle, etag, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(same)

	if d.Rev() != same.Rev() {
		t.Fatalf("Documnt with handle %s not fetched: %s != %s", handle, d.Id(), same.Rev())
	}

	log.Printf("%s == %s", d.Rev(), same.Rev())
}
