package main

import (
	"bytes"
	"encoding/json"
    "database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func retrieveMessages(id string, db *sql.DB) ([]Message, error) {
	query := `SELECT role, message FROM conversations WHERE conversationId = ? ORDER BY timestamp ASC`

	// Execute the query
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an array of Messages
	var messages []Message

	// Iterate over the rows and extract role and message
	for rows.Next() {
		var role, message string
		if err := rows.Scan(&role, &message); err != nil {
			return nil, err
		}
		messages = append(messages, Message{Role: role, Content: message})
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

    fmt.Println("Dumping contents from retrieveMessages:")
    for _, msg := range messages {
        fmt.Printf("Role: %s, Content: %s\n", msg.Role, msg.Content)
    }

	// Return the retrieved messages
	return messages, nil
}

func storeMessage(data Request, db *sql.DB) (bool, error) {
    // Prepare the SQL statement
	stmt, err := db.Prepare(`INSERT INTO conversations (message, conversationId, role) VALUES (?, ?, ?);`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL statement with the response values
	_, err = stmt.Exec(data.Message, data.Id, data.Role)
	if err != nil {
        log.Printf("Error storing message: %v", err)
        return false, err
	}

    log.Printf("Logged new conversation: %v", data)
	return true, nil
}

func getCompletionFromOpenAi(prompt Query) (Request, error) {
	apiEndpoint := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("CHAT_GPT_API_KEY")

    // Encode the request body into JSON
	jsonBody, err := json.Marshal(prompt)
	if err != nil {
		errMsg := "Error encoding the request data: " + err.Error()
		log.Println(errMsg)
		return Request{}, err
	}
    fmt.Println("String-ified json data for Open AI API request: ", string(jsonBody))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		errMsg := "Error creating the HTTP request: " + err.Error()
		log.Println(errMsg)
		return Request{}, err
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Request{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API request failed with status code: %d", resp.StatusCode)
		if err != nil {
			errMsg += " " + err.Error()
		}
		log.Println(errMsg)
		return Request{}, fmt.Errorf("%d", resp.StatusCode)
	}

	// Read the response body into a byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := "Error reading the response body from Chat GPT API: " + err.Error()
		log.Println(errMsg)
		return Request{}, err
	}
    // Log API response
    log.Print(string(body))

	// Unmarshal the JSON into the Response struct
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return Request{}, err
	}

	var lastMessageContent string
	lastChoiceIndex := len(response.Choices) - 1
	if lastChoiceIndex >= 0 {
		lastMessageContent = response.Choices[lastChoiceIndex].Message.Content
	} else {
		return Request{}, nil
	}

    request := Request{
        Id: response.Id,
        Message: lastMessageContent,
        Role: "assistant",
    }
	return request, nil
}

