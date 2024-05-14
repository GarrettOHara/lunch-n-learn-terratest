package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func handleConversationIdCache() string {
	content, err := os.ReadFile(".lastConversationId")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(".lastConversationId")
			if err != nil {
				log.Printf("Error creating file: %v", err)
				return ""
			}
			defer file.Close()
			return ""
		}
	}
	return string(content)
}

func cacheConversationId(body []byte) {
    file, err := os.OpenFile(".lastConversationId", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}
	var request Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		return
	}
    _, err = file.Write([]byte(request.Id))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}
}

func startNewConversation() ([]byte, error) {
    res, err := http.Get("http://localhost:8080/chatbot")
	if err != nil {
		log.Printf("Error sending initial GET request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

    // Convert to byte array
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}
    cacheConversationId(body)

	return body, nil
}

func pickupLastConversation(conversationId string) ([]byte, error) {
	// TODO: write this function to grab the last message
	// in the conversation and return it to print out to
	// the console
	//
	//
	// send simple API call to API with continueChatSession
	// and API token
	//
	// Server will send back last message and then invoke user to
	// response, continuing the conversation
	//
	// Chat logs are stored on the server side via sqlite3

	log.Printf("Conversation ID grapped from cache: %s", conversationId)
	return []byte{}, nil
}

func handleConversation(conversationId string) {
	// Handle user input until Ctrl+C is pressed
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			fmt.Println("\nExiting...")
			os.Exit(0)
		default:
			// Wait for user input
			fmt.Print("")
			reader := bufio.NewReader(os.Stdin)
			userInput, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading user input:", err)
				continue
			}

			// Send POST request with user input
			userInput = strings.TrimSpace(userInput)
			if userInput != "" {
				// Create a Request object
				requestData := Request{
					Id:      conversationId,
					Message: userInput,
                    Role:    "user",
				}
				// Marshal the Request object into JSON format
				jsonData, err := json.Marshal(requestData)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
					return
				}
				// Make API call
				resp, err := http.Post("http://localhost:8080/chatbot", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Printf("Error sending POST request: %v", err)
					continue
				}
				defer resp.Body.Close()

				// Read response
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading response body: %v", err)
					continue
				}
                fmt.Println("Response from POST reqeuest:")
                fmt.Println(string(body))

				// Format and print to CLI
				_ = handleRequest(body, false)
			}
		}
	}
	// Block indefinitely
	// select {}
}

func handleRequest(body []byte, initial bool) Request {
	var request Request
	err := json.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("Error unmarshalling request body PENIS:", err)
		return Request{}
	}
	if initial {
		fmt.Println("\n>", request.Message)
		fmt.Println("\n\nTo respond, type message, then hit Enter to send!")
		fmt.Println("To end the conversation, kill the process with '^C Enter'")
		fmt.Print("\n\n")
	} else {
		fmt.Print("\n> ", request.Message)
		fmt.Print("\n\n")
	}

	return request
}

func main() {
	var body []byte
	var err error
    body, err = startNewConversation()
    if err != nil {
        return
    }
	// conversationId := handleConversationIdCache()
	// if conversationId == "" {
	// 	body, err = startNewConversation()
	// 	if err != nil {
	// 		return
	// 	}
	// } else {
	// 	body, err = pickupLastConversation(conversationId)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	response := handleRequest(body, true)
	handleConversation(response.Id)
}
