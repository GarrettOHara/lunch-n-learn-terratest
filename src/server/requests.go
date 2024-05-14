package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, "OK\n"); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func queryChatBotGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	// Sample init for conversation
	// TOOD: add options
	text := `Pretend to be a super young zoomer/gen z 20 year old that speaks in sentences that barley make any sense, ask for their name and then have a conversation.`

	// Clean string
	// text = strings.ReplaceAll(text, "\n", "")
	// text = strings.ReplaceAll(text, "\t", "")
	// text = strings.TrimSpace(text)

	// Construct query
	var prompt Query
	prompt.Model = "gpt-3.5-turbo-0125"
	prompt.Messages = []Message{
		{Role: "user", Content: text},
	}
	fmt.Printf("%+v\n", prompt)

	// Send request to OpenAi API
	var request Request
	var err error
	request, err = getCompletionFromOpenAi(prompt)
	if err != nil {
		errMsg := "Error getting completion from API: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	save_response := request.Message
	request.Message = text
	request.Role = "user"
	// Store user prompt message in database
	_, err = storeMessage(request, db)
	if err != nil {
		errMsg := "Error saving chat gpt response to database: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	request.Message = save_response
	request.Role = "assistant"
	// Store chat gpt response message in database
	_, err = storeMessage(request, db)
	if err != nil {
		errMsg := "Error saving chat gpt response to database: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Handle/transform response data
	jsonData, err := json.Marshal(request)
	if err != nil {
		errMsg := "Error transforming the response data" + err.Error()
		log.Println(errMsg)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	// Response back to the client request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		errMsg := "Error responding back to client request: " + err.Error()
		log.Println(errMsg)
		return
	}
}

func queryChatBotPost(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	// Extract data from http request
	var data Request
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Clean string
	data.Message = strings.ReplaceAll(data.Message, "\n", "")
	data.Message = strings.ReplaceAll(data.Message, "\t", "")

	// Append user input to conversation list
	_, err = storeMessage(data, db)
	if err != nil {
		errMsg := "Error saving chat gpt response to database: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Grab conversation history from sqlite
	var messages []Message
	messages, err = retrieveMessages(data.Id, db)
	if err != nil {
		log.Printf("Error when retreiveing messages from database: %v", err)
		return
	}

	// Construct query
	var prompt Query
	prompt.Model = "gpt-3.5-turbo-0125"
	prompt.Messages = []Message{
		// {Role: "user"},
		// {Content: data.Message},
	}

	// Add all existing messages to prompt
	prompt.Messages = append(prompt.Messages, messages...)
	// appendMessages(&prompt, messages)

	// Send request to OpenAi API
	var request Request
	request, err = getCompletionFromOpenAi(prompt)
	if err != nil {
		errMsg := "Error getting completion from API: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Append chat gpt response to conversation list
	_, err = storeMessage(request, db)
	if err != nil {
		errMsg := "Error saving chat gpt response to database: " + err.Error()
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Handle/transform request data
	jsonData, err := json.Marshal(request)
	if err != nil {
		errMsg := "Error transforming the response data" + err.Error()
		log.Println(errMsg)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	// Response back to the client request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		errMsg := "Error responding back to client request: " + err.Error()
		log.Println(errMsg)
		return
	}
}
