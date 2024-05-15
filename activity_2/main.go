package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Quote struct {
	ID           string   `json:"_id"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	Tags         []string `json:"tags"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	DateAdded    string   `json:"dateAdded"`
	DateModified string   `json:"dateModified"`
}

func main() {
	url := "https://api.quotable.io/quotes/random"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var quotes []Quote
	err = json.NewDecoder(resp.Body).Decode(&quotes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(quotes) > 0 {
		quote := quotes[0]
		fmt.Printf("\n\n\"%s\"", quote.Content)
		fmt.Println(" -", quote.Author)
		fmt.Printf("\n\n")
	} else {
		fmt.Println("No quotes found in response")
	}
}
