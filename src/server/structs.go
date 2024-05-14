package main

type Response struct {
	Id      string   `json:"id"`
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

type Request struct {
	Id      string `json:"id"`
	Message string `json:"message"`
	Role    string `json:"role"`
}

type Query struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
