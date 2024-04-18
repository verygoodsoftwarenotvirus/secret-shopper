package main

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func fetchAllbirdsProducts(browser playwright.Browser) ([]*ProductGroup, error) {
	urls := []string{
		"https://www.allbirds.com/collections/mens",
		"https://www.allbirds.com/collections/womens",
		"https://www.allbirds.com/collections/little-kids",
		"https://www.allbirds.com/collections/big-kids",
		"https://www.allbirds.com/collections/socks",
	}

	products := []*ProductGroup{}

	for _, u := range urls {
		page, err := browser.NewPage()
		if err != nil {
			return nil, fmt.Errorf("could not create page: %w", err)
		}

		if _, err = page.Goto(u); err != nil {
			return nil, fmt.Errorf("could not goto: %w", err)
		}

		// todo: get all product links and then parse their pages
		allColorwayLinks, err := page.Locator(".Colorway__link").All()
		if err != nil {
			return nil, fmt.Errorf("could not find colorway links: %w", err)
		}

		allLinks := []string{}
		for colorwayIndex, colorway := range allColorwayLinks {
			colorwayLink, err := colorway.GetAttribute("href")
			if err != nil {
				slog.Error("could not get colorway link", slog.Any("error", err), slog.String("url", page.URL()), slog.Int("index", colorwayIndex))
				continue
			}

			if colorwayLink != "" {
				allLinks = append(allLinks, page.URL()+colorwayLink)
			}
		}

		for colorwayIndex, colorway := range allLinks {
			productGroup, err := parseAllbirdsProductPage(page, colorway)
			if err != nil {
				slog.Error("could not parse product", slog.Any("error", err), slog.String("url", page.URL()), slog.Int("index", colorwayIndex))
				return nil, fmt.Errorf("could not parse product group: %w", err)
			}

			products = append(products, productGroup)
		}
	}

	return products, nil
}

func parseAllbirdsProductPage(page playwright.Page, productURL string) (*ProductGroup, error) {
	slog.Info("fetching product data", slog.String("url", productURL))

	productGroup := &ProductGroup{
		URL:      productURL,
		Products: []*Product{},
	}

	if _, err := page.Goto(productURL); err != nil {
		return nil, fmt.Errorf("could not goto: %w", err)
	}

	time.Sleep(1500 * time.Millisecond)

	allColorways, err := page.Locator("button.ColorSwatchButton").All()
	if err != nil {
		return nil, fmt.Errorf("could not get colorways: %w", err)
	}

	if len(allColorways) == 0 {
		return nil, errors.New("empty colorways list")
	}

	for i, colorway := range allColorways {
		product := &Product{
			Attributes: map[string]string{},
		}

		// close any modals
		page.Locator(".CloseIcon").First().Click()

		clickTimeout := float64(5000)
		if err = colorway.Click(playwright.LocatorClickOptions{Timeout: &clickTimeout}); err != nil {
			return nil, fmt.Errorf("could not click colorway #%d: %w", i, err)
		}
		product.URL = page.URL()

		slog.Info("fetching product", slog.String("url", product.URL))

		productNameElement := page.Locator("h1.typography--secondary-heading").First()
		productGroup.Name, err = productNameElement.InnerText()
		if err != nil {
			return nil, fmt.Errorf("could not get product name: %w", err)
		}
		product.Name = productGroup.Name

		colorName, err := page.Locator(".Overview__name").First().InnerText()
		if err != nil {
			return nil, fmt.Errorf("could not get color name: %w", err)
		}
		product.Attributes["color"] = colorName

		imagesData, err := page.Locator(".ThumbnailButton > img").All()
		if err != nil {
			return nil, fmt.Errorf("could not get image thumbnails: %w", err)
		}

		for _, img := range imagesData {
			var rawSourceURL string
			rawSourceURL, err = img.GetAttribute("src")
			if err != nil {
				return nil, fmt.Errorf("could not get raw image source: %w", err)
			}

			rawSourceURL = strings.TrimPrefix(rawSourceURL, "https://cdn.allbirds.com/image/fetch/q_auto,f_auto/w_120,f_auto,q_auto/")
			product.ImageURLs = append(product.ImageURLs, rawSourceURL)
		}

		sizesData, err := page.Locator(".SizeButton").All()
		if err != nil {
			return nil, fmt.Errorf("could not get size data: %w", err)
		}

		for _, sizeContainer := range sizesData {
			size, err := sizeContainer.GetAttribute("aria-label")
			if err != nil {
				return nil, fmt.Errorf("could not get size aria label: %w", err)
			}

			size = strings.TrimPrefix(size, "Add Size ")
			product.Sizes = append(product.Sizes, size)
		}

		productGroup.Products = append(productGroup.Products, product)
	}

	return productGroup, nil
}
