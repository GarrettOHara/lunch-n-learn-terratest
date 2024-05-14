package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func openDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Create API log file
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Check if the database file exists
	if _, err := os.Stat("chat.db"); os.IsNotExist(err) {
		// Execute the bash script to create the database
		cmd := exec.Command("bash", "create_database.sh")
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Error executing script: %v\n", err)
		}
		log.Printf("%s", output)
	}

	// Create db connection
	db, err := openDatabase("chat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Setup API routing with a router
	router := http.NewServeMux()
	router.HandleFunc("/", healthCheck)
	router.HandleFunc("/chatbot", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			queryChatBotGet(w, r, db)
		} else if r.Method == http.MethodPost {
			queryChatBotPost(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	//TODO: http.HandleFunc("/getLastMessage", getLastMessage)

	// Start API on port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}
