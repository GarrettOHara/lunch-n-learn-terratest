
import (
	"bytes"
	"encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/", healthCheck)
    http.HandleFunc("/chatbot", queryChatBot)

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    if _, err := io.WriteString(w, "OK\n"); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func queryChatBot(w http.ResponseWriter, req *http.Request) {
	// Parse the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Extract the text from the request
	text := string(body)
    prompt := map[string]string{}

	// Check if the text is empty (beginning of conversation)
	if text == "" {
        text = "Pretend to be a super young zoomer/gen z 20 year old that speaks in sentences that barley make any sense, ask for their name and then have a conversation."
        prompt = map[string]string{"role": "system", "content": text}
    } else {
        prompt = map[string]string{"role": "user", "content": text}
    }

	// Make a request to the ChatGPT completion API
	id, response, err := getCompletionFromAPI(prompt)
	if err != nil {
        errorMessage := fmt.Sprintf("Error getting completion from API: %s", err.Error())
        http.Error(w, errorMessage, http.StatusInternalServerError)
        fmt.Println(errorMessage) // Log the error message
        return
	}

    // Create ResponseData struct
	responseData := ResponseData{
		ID:      id,
		Message: response,
	}

	// Encode response map as JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getCompletionFromAPI(prompt map[string]string) (string, string, error) {
	// Define the API endpoint
	apiEndpoint := "https://api.openai.com/v1/chat/completions"

	// Your API key
	apiKey := os.Getenv("CHAT_GPT_API_KEY")

    // Define the request body parameters
	requestBody := map[string]interface{}{
		"model":           "gpt-3.5-turbo-0125",
		"messages": []map[string]string{
			prompt,
		},
	}

	// Encode the request body into JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", err
	}

	// Set the request headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

    // Read the response body into a byte slice
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", err
    }

    // Unmarshal the JSON into the Response struct
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}

    // Extract the ID
	id := response.ID
    var lastMessageContent string
	lastChoiceIndex := len(response.Choices) - 1
	if lastChoiceIndex >= 0 {
		lastMessageContent = response.Choices[lastChoiceIndex].Message.Content
	} else {
        return "", "", nil
	}

    return id, lastMessageContent, nil
}
