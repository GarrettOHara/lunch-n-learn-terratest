package main

// Response represents the structure of the JSON response
type Response struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"` // Assuming it can be null
	FinishReason string      `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ResponseData represents the structure of the JSON response
type ResponseData struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// ChatGptRequest represents the request body parameters
type ChatGptRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}
