package main

// Response represents the structure of the JSON response
type Response struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

// Choice represents the structure of each choice in the choices list
// type Choice struct {
// 	Message struct {
// 		Content string `json:"content"`
// 	} `json:"message"`
// }

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
