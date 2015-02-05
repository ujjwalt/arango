package arango

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	endpoint string       // the endpoint to connect to
	client   *http.Client // the client used for connecting to the endpoint
)

const (
	apiChanedPanic = "The API has changed! Please open an issue immediately"
)

func Connect(rawurl string) {
	endpoint = rawurl
	client = &http.Client{}
}

func CurrentDB() (*Database, error) {
	obj, err := get("/_api/database/current")
	if err != nil {
		return nil, err
	}

	obj, ok := obj["result"].(map[string]interface{})
	if !ok {
		panic(apiChanedPanic)
	}

	d := &Database{
		Id:       obj["id"].(string),
		Name:     obj["name"].(string),
		Path:     obj["path"].(string),
		IsSystem: obj["isSystem"].(bool),
	}
	return d, nil
}

func get(path string) (o map[string]interface{}, err error) {
	data, err := getRaw(path)
	err = json.Unmarshal(data, &o)
	return
}

func getRaw(path string) (data []byte, err error) {
	r, err := client.Get(endpoint + path)
	if err != nil {
		return
	}
	defer r.Body.Close()

	data, err = ioutil.ReadAll(r.Body)
	return
}
