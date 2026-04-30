package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"net/http"

	"sigs.k8s.io/yaml"
)

type ErrorResponse struct {
	Message  string   `json:"message"`
	Status   int      `json:"status"`
}

type CharacterListResponse struct {
	Message   string         `json:"message"`
	Status    int            `json:"status"`
	Characters []CharacterListEntry `json:"characters"`
}

type CharacterListEntry struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
}

func StartRestController(port string) {
	// Define a route and handler for serving static files at the root
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/skills", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			res := ErrorResponse {
				Message: fmt.Sprintf("Bad request. Method not supported."),
				Status: 400,
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		yamlData, _ := os.ReadFile("./data/skills.yaml")

		jsonData, err := yaml.YAMLToJSON(yamlData)
		if err != nil {
			http.Error(w, "Could not convert to JSON", http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})

	http.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			res := ErrorResponse {
				Message: fmt.Sprintf("Bad request. Method not supported."),
				Status: 400,
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		files, err := os.ReadDir(".") // Use "." for current directory or provide a path
		if err != nil {
			log.Fatal(err)
		}

		charList := make([]CharacterListEntry, 0, len(files))
		for _, entry := range files {
			entry := CharacterListEntry{
				Filename: entry.Name(),
				Title:    entry.Name(),
			}

			charList = append(charList, entry)
		}

		res := CharacterListResponse {
			Message:   "Success",
			Status:     200,
			Characters: charList,
		}
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/characters/{name}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			res := ErrorResponse {
				Message: fmt.Sprintf("Bad request. Method not supported."),
				Status: 400,
			}
			json.NewEncoder(w).Encode(res)
			return
		}

		charName := r.PathValue("name")
		charPath := fmt.Sprintf("./data/characters/%s", charName)

		yamlData, _ := os.ReadFile(charPath)

		jsonData, err := yaml.YAMLToJSON(yamlData)
		if err != nil {
			http.Error(w, "Could not convert to JSON", http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})

	// Start the server on port 8080
	log.Println("Server starting on http://localhost:"+port+"...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
