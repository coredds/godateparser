//go:build examples
// +build examples

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/coredds/godateparser"
)

// CLI Tool Example demonstrates a command-line date parser

func main() {
	// Define flags
	dateOrder := flag.String("order", "MDY", "Date component order: YMD, MDY, or DMY")
	languages := flag.String("lang", "en", "Comma-separated language codes (e.g., en,es,fr)")
	preferFuture := flag.Bool("future", false, "Prefer future dates for ambiguous inputs")
	format := flag.String("format", "2006-01-02", "Output date format (Go time format)")
	extract := flag.Bool("extract", false, "Extract all dates from input text")
	relative := flag.String("relative", "", "Base date for relative parsing (YYYY-MM-DD)")
	verbose := flag.Bool("v", false, "Verbose output")
	interactive := flag.Bool("i", false, "Interactive mode")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [date_string]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A command-line date parser supporting multiple formats and languages.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s \"December 31, 2024\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -order=DMY \"31/12/2024\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -lang=es \"31 diciembre 2024\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -extract \"Meeting on Dec 31 and Jan 15\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -future \"next Monday\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSupported languages: en, es, pt, fr, de, it, nl, ru, zh, ja\n")
	}

	flag.Parse()

	// Interactive mode
	if *interactive {
		runInteractiveMode(*dateOrder, *languages, *preferFuture, *format, *relative, *verbose)
		return
	}

	// Get input from arguments or stdin
	var input string
	if flag.NArg() > 0 {
		input = strings.Join(flag.Args(), " ")
	} else {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}
		if input == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	// Build settings
	settings := buildSettings(*dateOrder, *languages, *preferFuture, *relative)

	if *verbose {
		fmt.Fprintf(os.Stderr, "Input: %s\n", input)
		fmt.Fprintf(os.Stderr, "Settings: DateOrder=%s, Languages=%v\n",
			settings.DateOrder, settings.Languages)
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Extract or parse
	if *extract {
		extractDates(input, settings, *format, *verbose)
	} else {
		parseSingleDate(input, settings, *format, *verbose)
	}
}

func buildSettings(dateOrder, languages string, preferFuture bool, relative string) *godateparser.Settings {
	settings := &godateparser.Settings{
		DateOrder:    dateOrder,
		Languages:    strings.Split(languages, ","),
		RelativeBase: time.Now(),
	}

	if preferFuture {
		settings.PreferDatesFrom = "future"
	}

	if relative != "" {
		baseDate, err := time.Parse("2006-01-02", relative)
		if err == nil {
			settings.RelativeBase = baseDate
		}
	}

	return settings
}

func parseSingleDate(input string, settings *godateparser.Settings, format string, verbose bool) {
	parsed, err := godateparser.ParseDate(input, settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Parsed successfully!\n")
		fmt.Printf("  ISO 8601:  %s\n", parsed.Format(time.RFC3339))
		fmt.Printf("  Unix Time: %d\n", parsed.Unix())
		fmt.Printf("  Custom:    %s\n", parsed.Format(format))
		fmt.Printf("  Human:     %s\n", parsed.Format("Monday, January 2, 2006 at 3:04 PM"))
	} else {
		fmt.Println(parsed.Format(format))
	}
}

func extractDates(input string, settings *godateparser.Settings, format string, verbose bool) {
	dates, err := godateparser.ExtractDates(input, settings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(dates) == 0 {
		if verbose {
			fmt.Fprintln(os.Stderr, "No dates found in input")
		}
		return
	}

	if verbose {
		fmt.Printf("Found %d date(s):\n\n", len(dates))
		for i, d := range dates {
			fmt.Printf("%d. \"%s\" (position %d, confidence %.2f)\n",
				i+1, d.MatchedText, d.Position, d.Confidence)
			fmt.Printf("   → %s\n\n", d.Date.Format(format))
		}
	} else {
		for _, d := range dates {
			fmt.Println(d.Date.Format(format))
		}
	}
}

func runInteractiveMode(dateOrder, languages string, preferFuture bool, format, relative string, verbose bool) {
	fmt.Println("=== Interactive Date Parser ===")
	fmt.Println("Enter dates to parse (or 'help' for commands, 'quit' to exit)")
	fmt.Println()

	settings := buildSettings(dateOrder, languages, preferFuture, relative)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		switch strings.ToLower(input) {
		case "quit", "exit", "q":
			fmt.Println("Goodbye!")
			return

		case "help", "h", "?":
			printHelp()
			continue

		case "settings":
			fmt.Printf("Current settings:\n")
			fmt.Printf("  Date Order:    %s\n", settings.DateOrder)
			fmt.Printf("  Languages:     %v\n", settings.Languages)
			fmt.Printf("  Prefer Future: %v\n", settings.PreferDatesFrom == "future")
			fmt.Printf("  Format:        %s\n", format)
			fmt.Printf("  Relative Base: %s\n", settings.RelativeBase.Format("2006-01-02"))
			continue

		case "examples":
			printExamples()
			continue
		}

		// Try to parse
		parsed, err := godateparser.ParseDate(input, settings)
		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		// Show result
		fmt.Printf("✓ Parsed: %s\n", parsed.Format(format))
		if verbose {
			fmt.Printf("  Full: %s\n", parsed.Format(time.RFC3339))
			fmt.Printf("  Unix: %d\n", parsed.Unix())

			// Calculate relative time
			now := time.Now()
			diff := parsed.Sub(now)
			if diff > 0 {
				fmt.Printf("  In:   %s\n", formatDuration(diff))
			} else {
				fmt.Printf("  Ago:  %s\n", formatDuration(-diff))
			}
		}
		fmt.Println()
	}
}

func printHelp() {
	fmt.Println("\nCommands:")
	fmt.Println("  help      - Show this help")
	fmt.Println("  settings  - Show current settings")
	fmt.Println("  examples  - Show example inputs")
	fmt.Println("  quit      - Exit interactive mode")
	fmt.Println("\nJust type any date string to parse it!")
	fmt.Println()
}

func printExamples() {
	fmt.Println("\nExample date inputs:")
	fmt.Println("  Absolute:")
	fmt.Println("    2024-12-31")
	fmt.Println("    December 31, 2024")
	fmt.Println("    31/12/2024")
	fmt.Println()
	fmt.Println("  Relative:")
	fmt.Println("    yesterday")
	fmt.Println("    next Monday")
	fmt.Println("    in 2 weeks")
	fmt.Println("    3 days ago")
	fmt.Println()
	fmt.Println("  Time:")
	fmt.Println("    3:30 PM")
	fmt.Println("    quarter past 3")
	fmt.Println("    noon")
	fmt.Println()
	fmt.Println("  Other:")
	fmt.Println("    Q4 2024")
	fmt.Println("    2024-W15")
	fmt.Println("    1609459200")
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}

/*
Usage Examples:

# Parse a simple date
./cli_tool "December 31, 2024"

# Parse with custom date order
./cli_tool -order=DMY "31/12/2024"

# Parse in Spanish
./cli_tool -lang=es "31 diciembre 2024"

# Parse relative date with future preference
./cli_tool -future "Monday"

# Extract all dates from text
./cli_tool -extract "Meeting on Dec 31 and follow-up on Jan 15"

# Custom output format
./cli_tool -format="January 2, 2006" "2024-12-31"

# Verbose output
./cli_tool -v "next Monday"

# Parse with custom relative base
./cli_tool -relative=2024-01-01 "in 2 weeks"

# Interactive mode
./cli_tool -i

# Read from stdin
echo "December 31, 2024" | ./cli_tool

# Process multiple dates from file
cat dates.txt | while read line; do ./cli_tool "$line"; done

# Use in shell scripts
DATE=$(./cli_tool "tomorrow")
echo "Tomorrow is: $DATE"

# Combined options
./cli_tool -lang=es,en -extract -v "Event: 31 diciembre or December 31"

# Multiple languages
./cli_tool -lang="fr,de,it" "31 décembre 2024"
*/
