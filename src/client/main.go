package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Send initial GET request
	res, err := http.Get("http://localhost:8080/chatbot")
	if err != nil {
		fmt.Println("Error sending initial GET request:", err)
		return
	}
	defer res.Body.Close()

	// Read and print response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}


	var response Response
    // Unmarshal response body into a Response struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}
    fmt.Println("\n>", response.Message)
    fmt.Println("\n\nTo respond, type message, then hit Enter to send!")
    fmt.Println("To end the conversation, kill the process with '^C Enter'")
    fmt.Print("\n\n")

	// Handle user input until Ctrl+C is pressed
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
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
					resp, err := http.Post("http://localhost:8080/chatbot", "text/plain", bytes.NewBufferString(userInput))
					if err != nil {
						fmt.Println("Error sending POST request:", err)
						continue
					}
					defer resp.Body.Close()

					// Read and print response
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("Error reading response body:", err)
						continue
					}

					// Unmarshal response body into a Response struct
					err = json.Unmarshal(body, &response)
					if err != nil {
						fmt.Println("Error unmarshalling response body:", err)
						continue
					}

					// Print the message from the response
					fmt.Print("\n> ", response.Message)
                    fmt.Print("\n\n")
				}
			}
		}
	}()

	// Block indefinitely
	select {}
}

