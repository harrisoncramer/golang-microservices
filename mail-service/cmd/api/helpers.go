package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

/* Parses JSON from an HTTP request into the third argument, data */
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 /* 1MB */

	/* Read up to 1MB from a request body into body */
	body := http.MaxBytesReader(w, r.Body, int64(maxBytes))

	/* Create a JSON decoder from the 1MB body reader */
	decoder := json.NewDecoder(body)

	/* Call the Decode method, converting 1MB stream into Golang structs */
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	/* Cannot have multiple JSON values */
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only a single JSON value")
	}

	return nil

}

/* Writes JSON data, status, and headers to an http.ResponseWriter */
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	/* Marshal the JSON data (convert to string) */
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	/* Set the headers */
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	/* Write the JSON to the ResponseWriter */
	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

/*
	Writes an error and error code (defaults to 400) to the ResponseWriter

and sets an error field of "true" to the JSON body response
*/
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error() /* The string of the error */

	/* Use our utility function from above */
	return app.writeJSON(w, statusCode, payload)
}
