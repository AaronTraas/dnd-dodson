package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"net/http"

	"sigs.k8s.io/yaml"
)

type CharacterListResponse struct {
	Message   string         `json:"message"`
	Status    int            `json:"status"`
	Characters []CharacterListEntry `json:"characters"`
}

type CharacterListEntry struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
}

type GameRules struct {
	Stats   interface{} `json:"stats"`
	Skills  interface{} `json:"skills"`
	Classes interface{} `json:"classes"`
}

func StartRestController(port string) {
	// Define a route and handler for serving static files at the root
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)

		if r.Method != "GET" {
			http.Error(w, "Bad request. Method not supported.", http.StatusBadRequest)
			return
		}

		skillsYamlData, _ := os.ReadFile("./data/rules/skills.yaml")
		var rawSkills interface{}
		if err := yaml.Unmarshal(skillsYamlData, &rawSkills); err != nil {
			http.Error(w, "Could not unmarshal YAML", http.StatusInternalServerError)
			return
		}

		statsYamlData, _ := os.ReadFile("./data/rules/stats.yaml")
		var rawStats interface{}
		if err := yaml.Unmarshal(statsYamlData, &rawStats); err != nil {
			http.Error(w, "Could not unmarshal YAML", http.StatusInternalServerError)
			return
		}

		classesYamlData, _ := os.ReadFile("./data/rules/classes.yaml")
		var rawClasses interface{}
		if err := yaml.Unmarshal(classesYamlData, &rawClasses); err != nil {
			http.Error(w, "Could not unmarshal YAML", http.StatusInternalServerError)
			return
		}

		res := GameRules {
			Stats: rawStats,
			Skills: rawSkills,
			Classes: rawClasses,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "GET" {
			http.Error(w, "Bad request. Method not supported.", http.StatusBadRequest)
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
			http.Error(w, "Bad request. Method not supported.", http.StatusBadRequest)
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
