package arango

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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

func get(path string) (v map[string]interface{}, err error) {
	data, err := getRaw(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
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

func post(path string, body map[string]interface{}) (v map[string]interface{}, err error) {
	data, err := postRaw(path, body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &v)
	return
}

func postRaw(path string, object map[string]interface{}) (data []byte, err error) {
	data, err = json.Marshal(object)
	if err != nil {
		return
	}

	body := strings.NewReader(string(data))
	r, err := client.Post(endpoint+path, "application/json", body)
	if err != nil {
		return
	}
	defer r.Body.Close()

	data, err = ioutil.ReadAll(r.Body)
	return
}
