package arango

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Represents a connection to a database
type Connection struct {
	url    string       // url for the connection
	client *http.Client // a client dedicated for this connection

	// Requests
	currDB        *http.Request // cached request for fetching the current database
	accessibleDBs *http.Request // cached request for fetching accessible databases
	dbs           *http.Request // cached request for fetching all databases
}

var DefaultConnection *Connection

// Create and return a new connection for the raw url passed
func NewConnection(rawurl string) *Connection {
	// Create a HTTP client pertaining to this connection
	client := &http.Client{}
	c := &Connection{url: rawurl, client: client}
	// Set as default connection if one does not exist
	if DefaultConnection == nil {
		DefaultConnection = c
	}
	return c
}

// Make the passed request and return the result
func (c *Connection) makeRequest(request *http.Request) ([]byte, error) {
	// Make the request
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Now read the body and parse it into a hash
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var h hash
	if err = json.Unmarshal(body, &h); err != nil {
		return nil, err
	}

	// Convert the result to a database
	// Marshall the result into a slice of bytes
	data, err := json.Marshal(h["result"])
	if err != nil {
		return nil, err
	}
	return data, nil
}

// A shorthand for writing a map that represents JSON objects
type hash map[string]interface{}

// Fetch the current database for the connection c
// GET /_api/database/current
func (c *Connection) CurrentDatabase() (*Database, error) {
	// Create the correct request and cache it if it is not cached
	if c.currDB == nil {
		req, err := http.NewRequest("GET", c.url+"/_api/database/current", nil)
		if err != nil {
			return nil, err
		}
		c.currDB = req
	}

	data, err := c.makeRequest(c.currDB)
	if err != nil {
		return nil, err
	}

	db := c.NewDatabase() // Initialize database with the current connection
	// Unmarshall the slice of bytes into the struct
	if err = json.Unmarshal(data, db); err != nil {
		return nil, err
	}
	DefaultDatabase = db
	return db, nil
}

func CurrentDatabase() (*Database, error) {
	return DefaultConnection.CurrentDatabase()
}

// List of accessible databases
// Retrieves the list of all databases the current user can access without specifying a different username or password.
// GET /_api/database/user
func (c *Connection) AccessibleDatabases() ([]string, error) {
	if c.accessibleDBs == nil {
		req, err := http.NewRequest("GET", c.url+"/_api/database/user", nil)
		if err != nil {
			return nil, err
		}
		c.accessibleDBs = req
	}

	data, err := c.makeRequest(c.accessibleDBs)
	if err != nil {
		return nil, err
	}

	var dbs []string
	// Unmarshall the slice of bytes into the struct
	if err = json.Unmarshal(data, &dbs); err != nil {
		return nil, err
	}

	return dbs, nil
}

func AccessibleDatabases() ([]string, error) {
	return DefaultConnection.AccessibleDatabases()
}

// List of databases
// Retrieves the list of all existing databases
// Note: retrieving the list of databases is only possible from within the _system database.
// GET /_api/database
func (c *Connection) Databases() ([]string, error) {
	if c.dbs == nil {
		req, err := http.NewRequest("GET", c.url+"/_api/database", nil)
		if err != nil {
			return nil, err
		}
		c.dbs = req
	}

	data, err := c.makeRequest(c.dbs)
	if err != nil {
		return nil, err
	}

	var dbs []string
	// Unmarshall the slice of bytes into the struct
	if err = json.Unmarshal(data, &dbs); err != nil {
		return nil, err
	}

	return dbs, nil
}

func Databases() ([]string, error) {
	return DefaultConnection.Databases()
}
