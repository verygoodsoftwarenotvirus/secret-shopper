package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type GoodOnYouPayload struct {
	Category []string         `json:"category"`
	Offset   int              `json:"offset"`
	Filters  GoodOnYouFilters `json:"filters"`
}

type GoodOnYouFilters struct {
	Category []string `json:"category"`
}

func fetchAllGoodOnYouData() error {
	categories := []string{"tops", "activewear", "bottoms", "accessories", "bags", "shoes", "plussize", "maternity", "swimwear", "sleepwear", "suits", "denim", "outerwear"}
	finalResults := []GoodOnYouBrand{}

	for _, category := range categories {
		results, err := fetchGoodOnYouResultForCategory(category)
		if err != nil {
			return err
		}

		for _, brand := range results.Result.Brands {
			finalResults = append(finalResults, brand)
		}
	}

	renderToJSONFile(finalResults, "good_on_you.json")

	return nil
}

func fetchGoodOnYouResultForCategory(category string) (*GoodOnYouQueryResultContainer, error) {
	data := GoodOnYouPayload{
		Category: []string{category},
		Offset:   10000,
		Filters: GoodOnYouFilters{
			Category: []string{category},
		},
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://public-api.goodonyou.eco/parse/functions/browseCategoriesV4", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("X-Parse-Application-Id", "gcrp2V42PHW7S8ElL639")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var x *GoodOnYouQueryResultContainer
	if err = json.NewDecoder(resp.Body).Decode(&x); err != nil {
		return nil, err
	}

	return x, nil
}

type GoodOnYouQueryResultContainer struct {
	Result GoodOnYouQueryResult `json:"result"`
}

type GoodOnYouQueryResult struct {
	Page       int                  `json:"page"`
	Total      int                  `json:"total"`
	Category   GoodOnYouCategory    `json:"category"`
	Categories []GoodOnYouCategory  `json:"categories"`
	Filters    GoodOnYouQueryFilter `json:"filters"`
	Picks      []GoodOnYouPick      `json:"picks"`
	Brands     []GoodOnYouBrand     `json:"brands"`
}

type GoodOnYouQueryFilter struct {
	Territories   []GoodOnYouTerritory     `json:"territories"`
	Locations     GoodOnYouLocations       `json:"locations"`
	Values        []GoodOnYouCompanyValues `json:"values"`
	Subcategories []GoodOnYouSubcategory   `json:"subcategories"`
}

type GoodOnYouCategory struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	Image         string                 `json:"image"`
	Gender        []string               `json:"gender"`
	Qty           int                    `json:"qty"`
	Total         GoodOnYouCategoryTotal `json:"total"`
	Subcategories []string               `json:"subcategories"`
}

type GoodOnYouCategoryTotal struct {
	All   int `json:"all"`
	Kids  int `json:"kids"`
	Men   int `json:"men"`
	Women int `json:"women"`
}

type GoodOnYouTerritory struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GoodOnYouLocations struct {
	Eu GoodOnYouRegion `json:"EU"`
	Na GoodOnYouRegion `json:"NA"`
	Cs GoodOnYouRegion `json:"CS"`
	Oc GoodOnYouRegion `json:"OC"`
	Af GoodOnYouRegion `json:"AF"`
	As GoodOnYouRegion `json:"AS"`
}

type GoodOnYouRegion struct {
	Code      string             `json:"code"`
	Name      string             `json:"name"`
	Countries []GoodOnYouCountry `json:"countries"`
}

type GoodOnYouCountry struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Continent string `json:"continent"`
}

type GoodOnYouCompanyValues struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GoodOnYouSubcategory struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Gender []string `json:"gender"`
}

type GoodOnYouPick struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Image         string `json:"image"`
	EthicalLabel  string `json:"ethical_label"`
	EthicalRating int    `json:"ethical_rating"`
	Price         int    `json:"price"`
	Order         int    `json:"order"`
}

type GoodOnYouBrand struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Image         string   `json:"image"`
	EthicalLabel  string   `json:"ethical_label"`
	EthicalRating int      `json:"ethical_rating"`
	Price         int      `json:"price"`
	Territory     []string `json:"territory"`
	Women         int      `json:"women"`
	Men           int      `json:"men"`
	Kids          int      `json:"kids"`
	IsStaffPick   bool     `json:"isStaffPick,omitempty"`
}
