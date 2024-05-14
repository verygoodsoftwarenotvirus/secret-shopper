package main

import (
	"encoding/json"
	"log"
	"os"
)

func renderToJSONFile(content any, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Encode the struct as JSON and write it to the file
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(content); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fetchAllFairWearData()
}
