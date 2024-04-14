package main

import (
	"encoding/json"
	"github.com/playwright-community/playwright-go"
	"log"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := playwright.Install(); err != nil {
		log.Fatalf("could not install playwright: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	headless := false
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	productGroup, err := fetchAllbirdsProducts(browser)
	if err != nil {
		log.Fatalf("could not fetch Allbirds products: %v", err)
	}

	if err = json.NewEncoder(os.Stdout).Encode(productGroup); err != nil {
		log.Fatalf("could not encode product group: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
