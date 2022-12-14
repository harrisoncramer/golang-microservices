package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

/* Test handler for requests to the "/" route */
func (app *Config) BrokerTest(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

/* Handles all other requests */
func (app *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// DEPRECATED: app.logItem(w, requestPayload.Log)
		app.logEventViaRabbit(w, requestPayload.Log)
	case "mail":
		if os.Getenv("APP_ENV") != "production" {
			app.sendMail(w, requestPayload.Mail)
		} else {
			app.errorJSON(w, errors.New("Mail service not available"))
		}

	default:
		app.errorJSON(w, errors.New("Unknown action"))

	}
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, entry LogPayload) {
	err := app.pushToQueue(entry.Name, entry.Data)
	if err != nil {
		app.errorJSON(w, err)
	}

	response := jsonResponse{
		Error:   false,
		Message: "Logged via RabbitMQ",
	}

	app.writeJSON(w, http.StatusAccepted, response)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	err = emitter.Push(string(jsonData), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

/* DEPRECATED: Sends a message to the logger-service and writes success message to requester */
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

	/* Create some JSON to send to the logger microservice */
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	/* Construct the request */
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	/* On success, send a success message to the requesting service */
	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}

/* Authenticates client or rejects authentication request */
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	/* Create some JSON to send to the auth microservice */
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	/* Call the service (build the HTTP request) */
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	/* Make sure we get back the correct status code */
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling authentication service"))
		return
	}

	/* Decode the json and ensure no errors are passed */
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	/* Forward the valid json from the authentication service */
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)

}

/* Sends mail with given message to the recipient specified */
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceUrl := "http://mail-service/send"

	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	/* Decode the json and ensure no errors are passed */
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling mail service"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Message sent to " + msg.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
