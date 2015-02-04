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
	currDB *http.Request // cached request for fetching the current database
}

var DefaultConnection *Connection

// Create and return a new connection for the raw url passed
func NewConnection(rawurl string) (*Connection, error) {
	// Create a HTTP client pertaining to this connection
	client := &http.Client{}
	c := &Connection{url: rawurl, client: client}
	// Set as default connection if one does not exist
	if DefaultConnection == nil {
		DefaultConnection = c
	}
	return c, nil
}

// Make the passed request and return the result
func (c *Connection) makeRequest(request *http.Request) (interface{}, error) {
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

	result := h["result"]
	return result, nil
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

	result, err := c.makeRequest(c.currDB)
	if err != nil {
		return nil, err
	}

	// Convert the result to a database
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	db := &Database{conn: c} // Initialize database with the current connection
	if err = json.Unmarshal(data, db); err != nil {
		return nil, err
	}

	return db, nil
}

func CurrentDatabase() (*Database, error) {
	return DefaultConnection.CurrentDatabase()
}
