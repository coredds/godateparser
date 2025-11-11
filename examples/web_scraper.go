//go:build examples
// +build examples

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coredds/godateparser"
)

// WebScraperExample demonstrates extracting dates from scraped web content
func main() {
	fmt.Println("=== Web Scraper Example ===")
	fmt.Println()

	// Simulate scraping a webpage (in real scenario, use http.Get())
	htmlContent := `
	<html>
		<body>
			<h1>Upcoming Events</h1>
			<div class="event">
				<h2>Tech Conference 2024</h2>
				<p>Date: December 15, 2024</p>
				<p>Registration closes: 12/01/2024</p>
			</div>
			<div class="event">
				<h2>Product Launch</h2>
				<p>Scheduled for next Monday at 2 PM</p>
				<p>Follow-up meeting: 2 weeks from launch</p>
			</div>
			<div class="event">
				<h2>Annual Report Release</h2>
				<p>Q4 2024 results available on 2025-01-15</p>
				<p>Investors meeting: January 20, 2025</p>
			</div>
		</body>
	</html>
	`

	// In a real scenario, you'd fetch the content:
	// resp, err := http.Get("https://example.com/events")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// htmlContent = string(body)

	// Extract all dates from the content
	settings := &godateparser.Settings{
		RelativeBase: time.Now(),
		DateOrder:    "MDY", // US format common on web
	}

	dates, err := godateparser.ExtractDates(htmlContent, settings)
	if err != nil {
		fmt.Printf("Error extracting dates: %v\n", err)
		return
	}

	// Display extracted dates
	fmt.Printf("Found %d dates in webpage:\n\n", len(dates))

	for i, dateInfo := range dates {
		fmt.Printf("%d. \"%s\" (position %d)\n", i+1, dateInfo.MatchedText, dateInfo.Position)
		fmt.Printf("   Parsed as: %s\n", dateInfo.Date.Format("Monday, January 2, 2006 15:04 MST"))
		fmt.Printf("   Confidence: %.2f\n", dateInfo.Confidence)

		// Calculate days until event
		daysUntil := int(time.Until(dateInfo.Date).Hours() / 24)
		if daysUntil > 0 {
			fmt.Printf("   Days until: %d days\n", daysUntil)
		} else if daysUntil < 0 {
			fmt.Printf("   Days ago: %d days\n", -daysUntil)
		} else {
			fmt.Printf("   Today!\n")
		}
		fmt.Println()
	}

	// Group dates by month
	fmt.Println("=== Dates Grouped by Month ===")
	datesByMonth := make(map[string][]godateparser.ParsedDate)
	for _, d := range dates {
		key := d.Date.Format("2006-01")
		datesByMonth[key] = append(datesByMonth[key], d)
	}

	for month, monthDates := range datesByMonth {
		fmt.Printf("\n%s: %d date(s)\n", month, len(monthDates))
		for _, d := range monthDates {
			fmt.Printf("  - %s\n", d.Date.Format("Jan 2, 2006"))
		}
	}

	// Example: Filter events in next 30 days
	fmt.Println("\n=== Events in Next 30 Days ===")
	now := time.Now()
	thirtyDaysLater := now.AddDate(0, 0, 30)

	upcomingCount := 0
	for _, d := range dates {
		if d.Date.After(now) && d.Date.Before(thirtyDaysLater) {
			fmt.Printf("- %s (%s)\n", d.MatchedText, d.Date.Format("Jan 2, 2006"))
			upcomingCount++
		}
	}

	if upcomingCount == 0 {
		fmt.Println("No events in the next 30 days")
	}

	// Demonstrate scraping with multiple languages
	fmt.Println("\n=== Multi-Language Web Scraping ===")

	multilingualContent := `
		Event Schedule:
		- English: December 31, 2024
		- Spanish: 31 diciembre 2024
		- French: 31 décembre 2024
		- German: 31 Dezember 2024
	`

	multiLangSettings := &godateparser.Settings{
		Languages: []string{"en", "es", "fr", "de"},
	}

	multiDates, err := godateparser.ExtractDates(multilingualContent, multiLangSettings)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Extracted %d dates from multilingual content:\n", len(multiDates))
	for _, d := range multiDates {
		fmt.Printf("- \"%s\" → %s\n", d.MatchedText, d.Date.Format("2006-01-02"))
	}

	// Example helper function for real HTTP scraping
	fmt.Println("\n=== Usage with Real HTTP Requests ===")
	fmt.Println("Example code:")
	fmt.Println(`
	func scrapeDates(url string) ([]godateparser.ParsedDate, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return godateparser.ExtractDates(string(body), nil)
	}
	`)

	fmt.Println("\n=== Web Scraper Example Complete ===")
}

// Helper function (for demonstration) - not used in simulated example above
func scrapeDates(url string) ([]godateparser.ParsedDate, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Remove HTML tags for better extraction
	content := removeHTMLTags(string(body))

	return godateparser.ExtractDates(content, nil)
}

// Simple HTML tag remover
func removeHTMLTags(html string) string {
	// Very basic tag removal - in production use golang.org/x/net/html
	result := html
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + " " + result[start+end+1:]
	}
	return result
}
