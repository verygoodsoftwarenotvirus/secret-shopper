package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"log/slog"
	"os"
	"strings"
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

	var b bytes.Buffer
	if err = json.NewEncoder(&b).Encode(productGroup); err != nil {
		log.Fatalf("could not encode product group: %v", err)
	}

	fmt.Println(strings.Repeat("\n", 5))
	fmt.Println(b.String())
	fmt.Println(strings.Repeat("\n", 5))

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
