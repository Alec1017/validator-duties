package duties

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Performs a GET request and unmarshals the JSON with the given target type
func GetRequest(endpoint string, target interface{}) error {
	// make the request
	resp, err := http.Get(ConsensusClientUrl + endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return json.Unmarshal(body, target)
}

func PostRequest(endpoint string, target interface{}, data interface{}) error {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Create the request
	req, err := http.NewRequest("POST", ConsensusClientUrl+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return json.Unmarshal(body, target)
}
