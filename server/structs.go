package main

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

type ResponseData struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ChatGptRequest struct {
	Model    string              `json:"model"`
	Messages []map[string]string `json:"messages"`
}
