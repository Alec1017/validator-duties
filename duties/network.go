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
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

// Performs a POST request and unmarshals the JSON wwith the given target type
func PostRequest(endpoint string, target interface{}, data interface{}) error {
	// Marshal the struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create the request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}
