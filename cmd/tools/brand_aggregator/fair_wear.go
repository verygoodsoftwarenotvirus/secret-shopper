package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Payload struct {
	Requests []Requests `json:"requests"`
}

type Params struct {
	Facets           []any  `json:"facets"`
	HighlightPostTag string `json:"highlightPostTag"`
	HighlightPreTag  string `json:"highlightPreTag"`
	Page             int    `json:"page"`
	Query            string `json:"query"`
	TagFilters       string `json:"tagFilters"`
}

type Requests struct {
	IndexName string `json:"indexName"`
	Params    Params `json:"params"`
}

func fetchFairWearPageOfResults(pageNumber int) (*FairWearSearchResults, error) {
	log.Printf("fetching pageNumber #%d\n", pageNumber)

	data := Payload{
		Requests: []Requests{
			{
				IndexName: "brands",
				Params: Params{
					Page: pageNumber,
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://fwf.my.site.com/public/services/apexrest/api/v1.0/search/algolia", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.fairwear.org")
	req.Header.Set("Dnt", "1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read the unzipped content
	unzippedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var x *FairWearSearchResults
	if err = json.NewDecoder(bytes.NewReader(unzippedData)).Decode(&x); err != nil {
		return nil, err
	}

	return x, nil
}

func fetchAllFairWearData() {
	finalResults := []FairWearSearchHit{}

	var (
		pageNumber int
	)
	for {
		results, err := fetchFairWearPageOfResults(pageNumber)
		if err != nil {
			log.Fatal(err)
		}

		for _, result := range results.Results {
			finalResults = append(finalResults, result.Hits...)
		}

		if len(results.Results) == 1 && len(results.Results[0].Hits) == 0 {
			break
		}

		pageNumber++
	}

	renderToJSONFile(finalResults, "fair_wear.json")
}

type FairWearSearchResults struct {
	Results []FairWearSearchResult `json:"results"`
}

type FairWearSearchResult struct {
	ServerTimeMS     any                 `json:"serverTimeMS"`
	Query            string              `json:"query"`
	ProcessingTimeMS any                 `json:"processingTimeMS"`
	Params           any                 `json:"params"`
	Page             any                 `json:"page"`
	NbPages          int                 `json:"nbPages"`
	NbHits           int                 `json:"nbHits"`
	Index            string              `json:"index"`
	HitsPerPage      int                 `json:"hitsPerPage"`
	Hits             []FairWearSearchHit `json:"hits"`
	ExhaustiveTypo   any                 `json:"exhaustiveTypo"`
	ExhaustiveNbHits any                 `json:"exhaustiveNbHits"`
	Exhaustive       any                 `json:"exhaustive"`
}

type FairWearSearchHit struct {
	Website       string          `json:"website"`
	Slug          string          `json:"slug"`
	ProductTypes  []ProductTypes  `json:"product_types"`
	Name          string          `json:"name"`
	MemberUntil   any             `json:"member_until"`
	MemberSince   int             `json:"member_since"`
	LogoSizes     LogoSizes       `json:"logo_sizes"`
	LogoPath      any             `json:"logo_path"`
	GalleryImages []any           `json:"gallery_images"`
	Description   string          `json:"description"`
	Company       FairWearCompany `json:"company"`
}

type LogoSizes struct {
	BrandLogo string `json:"brand-logo"`
}

type ProductTypes struct {
	ProductType string `json:"product_type"`
}

type FairWearCompany struct {
	TransparencyPercentage   string               `json:"transparency_percentage"`
	SourcingCountries        []SourcingCountry    `json:"sourcing_countries"`
	Rating                   string               `json:"rating"`
	NotablePractices         any                  `json:"notable_practices"`
	MemberUntil              any                  `json:"member_until"`
	MemberSince              int                  `json:"member_since"`
	LatestPerformanceCheck   any                  `json:"latest_performance_check"`
	HeadquartersLocationName any                  `json:"headquarters_location_name"`
	HeadquartersLocation     HeadquartersLocation `json:"headquarters_location"`
	HeaderSizes              HeaderSizes          `json:"header_sizes"`
	HeaderPath               string               `json:"header_path"`
	BenchmarkingScore        string               `json:"benchmarking_score"`
}

type SourcingCountry struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type HeadquartersLocation struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type HeaderSizes struct {
	HeroLarge string `json:"hero-large"`
	Hero      string `json:"hero"`
}
