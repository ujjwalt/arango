package arango

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Document struct {
	*hash
	_id  string        // document handle
	_key string        // document unique key
	_rev string        // document revision/etag
	req  *http.Request // cached request for this document
}

// A document handle uniquely identifies a document in the database.
// It is a string and consists of the collection's name and the document key (_key attribute) separated by /
func (d *Document) Id() string {
	return d._id
}

/*
A document key uniquely identifies a document in the collection it it stored in.
It can and should be used by clients when specific documents are searched.
Document keys are stored in the _key attribute of documents. The key values are
automatically indexed by ArangoDB in a collection's primary index.
Thus looking up a document by its key is regularly a fast operation. The _key
value of a document is immutable once the document has been created.
By default, ArangoDB will auto-generate a document key if no _key attribute is
specified, and use the user-specified _key otherwise.
This behavior can be changed on a per-collection level by creating collections
with the keyOptions attribute.
Using keyOptions it is possible to disallow user-specified keys completely, or
to force a specific regime for auto-generating the _key values.
*/
func (d *Document) Key() string {
	return d._key
}

/*
As ArangoDB supports MVCC, documents can exist in more than one revision. The
document revision is the MVCC token used to identify a particular revision of a
document. It is a string value currently containing an integer number and is
unique within the list of document revisions for a single document. Document
revisions can be used to conditionally update, replace or delete documents in
the database. In order to find a particular revision of a document, you need
the document handle and the document revision. ArangoDB currently uses 64bit
unsigned integer values to maintain document revisions internally. When
returning document revisions to clients, ArangoDB will put them into a string
to ensure the revision id is not clipped by clients that do not support big
integers. Clients should treat the revision id returned by ArangoDB as an
opaque string when they store or use it locally. This will allow ArangoDB to
change the format of revision ids later if this should be required. Clients
can use revisions ids to perform simple equality/non-equality comparisons
(e.g. to check whether a document has changed or not), but they should not use
revision ids to perform greater/less than comparisons with them to check if a
document revision is older than one another, even if this might work for some
cases.

Note: Revision ids have been returned as integers up to including ArangoDB 1.1
*/
func (d *Document) Rev() string {
	return d._rev
}

// Return the document with id in the database d
// GEt /_db/<database-name>/_api/document/<document-handle>
func (d *Database) Find(id string, v interface{}) error {
	url := d.conn.url + "/_db/" + d.Name + "/_api/document/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	data, err := d.conn.makeRequest(req)
	if err != nil {
		return err
	}

	// Unmarshall the slice of bytes into the struct
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// Return the document with id from the current database
// GET /_api/document/<document-handle>
func Find(id string, v interface{}) error {
	url := DefaultConnection.url + "/_api/document/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	data, err := DefaultConnection.makeRequest(req)
	if err != nil {
		return err
	}

	// Unmarshall the slice of bytes into the struct
	fmt.Println(string(data))
	if err = json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
