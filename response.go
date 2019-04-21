package main

import (
	"encoding/json"
	"net/http"
)

// Struct for send api data
type respAPI struct {
	Message interface{} `json:"message,omitempty"`
	Time    int64       `json:"time,omitempty"`
	Error   bool        `json:"error"`
}

// Create new reponse and send it 
func newResponse(writer http.ResponseWriter, response respAPI) error {
	b, err := json.Marshal(response)
	if err != nil {
		writer.Write([]byte("{ \"message\":  \"some error occurred\", \"error\": false}"))
		panic(err)
	}

	writer.Write(b)
	return nil
}
