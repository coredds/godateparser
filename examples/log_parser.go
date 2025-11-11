//go:build examples
// +build examples

package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/coredds/godateparser"
)

// LogParserExample demonstrates parsing dates from various log formats
func main() {
	fmt.Println("=== Log Parser Example ===")
	fmt.Println()

	// Simulate log files with different date formats
	logs := []string{
		"[2024-12-15 10:30:45] INFO: Application started",
		"[2024-12-15 10:31:02] DEBUG: Loading configuration",
		"[2024-12-15 10:31:15] ERROR: Connection failed to database",
		"ERROR: Failed to process request on December 15, 2024 at 2:30 PM",
		"WARNING: Memory usage high - logged yesterday at 15:00",
		"INFO: Scheduled maintenance next Monday",
		"CRITICAL: System reboot required in 2 hours",
		"[15/12/2024 10:45:30] Server response time: 250ms",
		"2024-12-15T10:50:00Z - API call successful",
		"Backup completed: 31 Dec 2024 23:59",
	}

	fmt.Println("=== Parsing Log Entries ===")

	baseTime := time.Date(2024, 12, 15, 12, 0, 0, 0, time.UTC)
	settings := &godateparser.Settings{
		RelativeBase: baseTime,
		DateOrder:    "YMD",
	}

	type LogEntry struct {
		Original  string
		Timestamp time.Time
		Message   string
	}

	var entries []LogEntry

	for _, log := range logs {
		dates, err := godateparser.ExtractDates(log, settings)
		if err != nil || len(dates) == 0 {
			fmt.Printf("Could not parse date from: %s\n", log)
			continue
		}

		// Use the first (most likely) date found
		entry := LogEntry{
			Original:  log,
			Timestamp: dates[0].Date,
			Message:   log,
		}
		entries = append(entries, entry)

		fmt.Printf("Log: %s\n", log)
		fmt.Printf("  → Timestamp: %s\n\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
	}

	// Analyze logs by time
	fmt.Println("=== Log Analysis ===")

	// Count logs per hour
	hourCounts := make(map[int]int)
	for _, entry := range entries {
		hourCounts[entry.Timestamp.Hour()]++
	}

	fmt.Println("\nLogs per hour:")
	for hour := 0; hour < 24; hour++ {
		if count := hourCounts[hour]; count > 0 {
			fmt.Printf("  %02d:00 - %2d log(s)\n", hour, count)
		}
	}

	// Find errors in last 24 hours
	fmt.Println("\n=== Recent Errors (last 24 hours) ===")
	cutoff := baseTime.Add(-24 * time.Hour)
	errorCount := 0

	for _, log := range logs {
		if strings.Contains(log, "ERROR") || strings.Contains(log, "CRITICAL") {
			dates, err := godateparser.ExtractDates(log, settings)
			if err == nil && len(dates) > 0 && dates[0].Date.After(cutoff) {
				fmt.Printf("  %s - %s\n", dates[0].Date.Format("15:04:05"), log)
				errorCount++
			}
		}
	}

	if errorCount == 0 {
		fmt.Println("  No recent errors")
	}

	// Parse different log formats
	fmt.Println("\n=== Common Log Formats ===")

	logFormats := map[string]string{
		"Apache Combined": `192.168.1.1 - - [15/Dec/2024:10:30:45 +0000] "GET /index.html HTTP/1.1" 200 1234`,
		"Nginx":           `2024/12/15 10:30:45 [error] 1234#0: *5678 connect() failed`,
		"Syslog":          `Dec 15 10:30:45 server kernel: [12345.678] Out of memory`,
		"JSON":            `{"timestamp":"2024-12-15T10:30:45Z","level":"info","message":"Request processed"}`,
		"Custom App":      `INFO  [2024-12-15 10:30:45.123] [main] Application starting`,
		"Windows Event":   `12/15/2024 10:30:45 AM - Application - Information - Event ID 1000`,
		"ISO 8601":        `2024-12-15T10:30:45.123456Z`,
		"Unix Timestamp":  `1702635045`,
		"Relative":        `Logged 2 hours ago`,
	}

	for format, logLine := range logFormats {
		dates, err := godateparser.ExtractDates(logLine, settings)
		if err == nil && len(dates) > 0 {
			fmt.Printf("%-20s → %s\n", format+":", dates[0].Date.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%-20s → (could not parse)\n", format+":")
		}
	}

	// Demonstrate streaming log processing
	fmt.Println("\n=== Streaming Log Processing ===")

	simulatedStream := `
[2024-12-15 10:00:00] Server started
[2024-12-15 10:05:23] First request received
[2024-12-15 10:10:45] Warning: High CPU usage
[2024-12-15 10:15:12] Request processed successfully
[2024-12-15 10:20:33] Error: Database timeout
`

	scanner := bufio.NewScanner(strings.NewReader(simulatedStream))
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if strings.TrimSpace(line) == "" {
			continue
		}

		dates, err := godateparser.ExtractDates(line, settings)
		if err == nil && len(dates) > 0 {
			// Process log entry
			timestamp := dates[0].Date

			// Detect log level
			level := "INFO"
			if strings.Contains(line, "Warning") {
				level = "WARN"
			} else if strings.Contains(line, "Error") {
				level = "ERROR"
			}

			fmt.Printf("Line %d [%s] %s - %s\n",
				lineNum, level, timestamp.Format("15:04:05"), line)
		}
	}

	// Example: Parse and aggregate logs by date
	fmt.Println("\n=== Daily Log Summary ===")

	dailyLogs := map[string][]string{
		"[2024-12-13 09:00:00] Backup started":    {},
		"[2024-12-13 23:00:00] Backup completed":  {},
		"[2024-12-14 10:00:00] System update":     {},
		"[2024-12-14 15:30:00] User login":        {},
		"[2024-12-15 08:00:00] Daily maintenance": {},
		"[2024-12-15 10:30:00] Error detected":    {},
	}

	logsByDay := make(map[string][]string)

	for log := range dailyLogs {
		dates, err := godateparser.ExtractDates(log, settings)
		if err == nil && len(dates) > 0 {
			dayKey := dates[0].Date.Format("2006-01-02")
			logsByDay[dayKey] = append(logsByDay[dayKey], log)
		}
	}

	for day, dayLogs := range logsByDay {
		fmt.Printf("\n%s (%d logs):\n", day, len(dayLogs))
		for _, log := range dayLogs {
			fmt.Printf("  - %s\n", log)
		}
	}

	// Helper function example
	fmt.Println("\n=== Usage Functions ===")
	fmt.Println(`
Example helper functions:

// Parse log timestamp
func parseLogTimestamp(logLine string) (time.Time, error) {
    dates, err := godateparser.ExtractDates(logLine, nil)
    if err != nil || len(dates) == 0 {
        return time.Time{}, errors.New("no date found")
    }
    return dates[0].Date, nil
}

// Filter logs by time range
func filterLogsByTime(logs []string, start, end time.Time) []string {
    var filtered []string
    for _, log := range logs {
        t, err := parseLogTimestamp(log)
        if err == nil && t.After(start) && t.Before(end) {
            filtered = append(filtered, log)
        }
    }
    return filtered
}

// Group logs by date
func groupLogsByDate(logs []string) map[string][]string {
    grouped := make(map[string][]string)
    for _, log := range logs {
        t, err := parseLogTimestamp(log)
        if err == nil {
            day := t.Format("2006-01-02")
            grouped[day] = append(grouped[day], log)
        }
    }
    return grouped
}
	`)

	fmt.Println("\n=== Log Parser Example Complete ===")
}
