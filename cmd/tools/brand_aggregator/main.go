package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type OnePercentForThePlanetResultsPage struct {
	Results      []OnePercentForThePlanetResults `json:"results"`
	TotalResults int                             `json:"totalResults"`
}

type OnePercentForThePlanetResults struct {
	ID       string                               `json:"id"`
	URI      string                               `json:"uri"`
	Name     string                               `json:"name"`
	Type     string                               `json:"type"`
	LogoURL  string                               `json:"logoUrl,omitempty"`
	Address  string                               `json:"address"`
	Snippet  string                               `json:"snippet,omitempty"`
	Location OnePercentForThePlanetResultLocation `json:"location"`
}

type OnePercentForThePlanetResultLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func fetch1PercentForThePlanetPage(pageNumber int) (*OnePercentForThePlanetResultsPage, error) {
	log.Printf("fetching page #%d\n", pageNumber)

	url := "https://e1k3unhdf2.execute-api.us-east-1.amazonaws.com/search?accountType=business&pageSize=20&"
	if pageNumber < 0 {
		url += fmt.Sprintf("page=%d", pageNumber)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://directories.onepercentfortheplanet.org")
	req.Header.Set("Dnt", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result *OnePercentForThePlanetResultsPage
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	finalResults := []OnePercentForThePlanetResults{}

	var (
		pageNumber, totalResults int
	)
	for {
		results, err := fetch1PercentForThePlanetPage(pageNumber)
		if err != nil {
			log.Fatal(err)
		}

		if totalResults == 0 {
			totalResults = results.TotalResults
		}

		finalResults = append(finalResults, results.Results...)

		if len(finalResults) >= totalResults {
			break
		}

		pageNumber++
	}

	// Open a file for writing
	file, err := os.Create("one_percent_for_the_planet.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Encode the struct as JSON and write it to the file
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(finalResults); err != nil {
		log.Fatal(err)
	}
}
