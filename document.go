package arango

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Document map[string]interface{}

func Find(handle string) (d Document, err error) {
	data, err := getRaw("/_api/document/" + handle)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &d)
	return
}

func FindIf(handle, etag string, match bool) (d Document, err error) {
	req, err := http.NewRequest("GET", endpoint+"/_api/document/"+handle, nil)
	if err != nil {
		return
	}
	if match {
		req.Header.Add("If-Match", etag)
	} else {
		req.Header.Add("If-None-Match", etag)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		// Parse the response into the document
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(data, d) // unmarshall into the document and return
	} else if !((match && resp.StatusCode == 412) || (!match && resp.StatusCode == 304) || resp.StatusCode == 404) {
		panic(apiChanedPanic) // did not get a expected status code based on match
	}
	return
}

func CreateDocument(doc map[string]interface{}, collection string, createCollection bool) (d Document, err error) {
	path := "/_api/document?collection=" + collection
	if createCollection {
		path += "&createCollection=true"
	}
	v, err := post(path, doc)
	if err != nil {
		return
	}
	d = Document(v)
	return
}
