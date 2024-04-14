package main

type ProductGroup struct {
	Name     string     `json:"name"`
	URL      string     `json:"url"`
	Products []*Product `json:"products"`
}

type Product struct {
	Name       string            `json:"name"`
	URL        string            `json:"url"`
	ImageURLs  []string          `json:"imageURLs"`
	Sizes      []string          `json:"sizes"`
	Attributes map[string]string `json:"attributes"`
}
