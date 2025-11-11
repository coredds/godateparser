package translations_test

import (
	"strings"
	"testing"
	"time"

	"github.com/coredds/godateparser/translations"
)

// Test Parse Month Helper

func TestParseMonth(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()
	french := translations.NewFrenchTranslation()

	tests := []struct {
		name      string
		input     string
		languages []*translations.Language
		wantMonth time.Month
		wantOK    bool
	}{
		{
			name:      "English full month",
			input:     "December",
			languages: []*translations.Language{english},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "English abbreviated month",
			input:     "Dec",
			languages: []*translations.Language{english},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "English lowercase",
			input:     "december",
			languages: []*translations.Language{english},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "Spanish month",
			input:     "diciembre",
			languages: []*translations.Language{spanish},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "French month",
			input:     "décembre",
			languages: []*translations.Language{french},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "Multiple languages - Spanish match",
			input:     "enero",
			languages: []*translations.Language{english, spanish},
			wantMonth: time.January,
			wantOK:    true,
		},
		{
			name:      "With whitespace",
			input:     "  December  ",
			languages: []*translations.Language{english},
			wantMonth: time.December,
			wantOK:    true,
		},
		{
			name:      "Invalid month",
			input:     "notamonth",
			languages: []*translations.Language{english},
			wantMonth: time.Month(0),
			wantOK:    false,
		},
		{
			name:      "Empty input",
			input:     "",
			languages: []*translations.Language{english},
			wantMonth: time.Month(0),
			wantOK:    false,
		},
		{
			name:      "No languages",
			input:     "December",
			languages: []*translations.Language{},
			wantMonth: time.Month(0),
			wantOK:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMonth, gotOK := translations.ParseMonth(tt.input, tt.languages...)

			if gotOK != tt.wantOK {
				t.Errorf("ParseMonth() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if gotMonth != tt.wantMonth {
				t.Errorf("ParseMonth() month = %v, want %v", gotMonth, tt.wantMonth)
			}
		})
	}
}

// Test Parse Weekday Helper

func TestParseWeekday(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()
	russian := translations.NewRussianTranslation()

	tests := []struct {
		name        string
		input       string
		languages   []*translations.Language
		wantWeekday time.Weekday
		wantOK      bool
	}{
		{
			name:        "English full weekday",
			input:       "Monday",
			languages:   []*translations.Language{english},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "English abbreviated weekday",
			input:       "Mon",
			languages:   []*translations.Language{english},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "English lowercase",
			input:       "monday",
			languages:   []*translations.Language{english},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "Spanish weekday",
			input:       "lunes",
			languages:   []*translations.Language{spanish},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "Russian weekday",
			input:       "понедельник",
			languages:   []*translations.Language{russian},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "Multiple languages - Spanish match",
			input:       "viernes",
			languages:   []*translations.Language{english, spanish},
			wantWeekday: time.Friday,
			wantOK:      true,
		},
		{
			name:        "With whitespace",
			input:       "  Monday  ",
			languages:   []*translations.Language{english},
			wantWeekday: time.Monday,
			wantOK:      true,
		},
		{
			name:        "Invalid weekday",
			input:       "notaweekday",
			languages:   []*translations.Language{english},
			wantWeekday: time.Weekday(0),
			wantOK:      false,
		},
		{
			name:        "Empty input",
			input:       "",
			languages:   []*translations.Language{english},
			wantWeekday: time.Weekday(0),
			wantOK:      false,
		},
		{
			name:        "No languages",
			input:       "Monday",
			languages:   []*translations.Language{},
			wantWeekday: time.Weekday(0),
			wantOK:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWeekday, gotOK := translations.ParseWeekday(tt.input, tt.languages...)

			if gotOK != tt.wantOK {
				t.Errorf("ParseWeekday() ok = %v, want %v", gotOK, tt.wantOK)
			}

			if gotWeekday != tt.wantWeekday {
				t.Errorf("ParseWeekday() weekday = %v, want %v", gotWeekday, tt.wantWeekday)
			}
		})
	}
}

// Test MatchesRelativeTerm

func TestMatchesRelativeTerm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		terms []string
		want  bool
	}{
		{
			name:  "Exact match",
			input: "yesterday",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  true,
		},
		{
			name:  "Case insensitive match",
			input: "Yesterday",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  true,
		},
		{
			name:  "With whitespace",
			input: "  yesterday  ",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  true,
		},
		{
			name:  "No match",
			input: "notfound",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  false,
		},
		{
			name:  "Empty input",
			input: "",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  false,
		},
		{
			name:  "Empty terms",
			input: "yesterday",
			terms: []string{},
			want:  false,
		},
		{
			name:  "Empty term in list",
			input: "yesterday",
			terms: []string{"", "yesterday", ""},
			want:  true,
		},
		{
			name:  "Partial match doesn't count",
			input: "yesterday morning",
			terms: []string{"yesterday", "today", "tomorrow"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translations.MatchesRelativeTerm(tt.input, tt.terms)
			if got != tt.want {
				t.Errorf("MatchesRelativeTerm(%q, %v) = %v, want %v",
					tt.input, tt.terms, got, tt.want)
			}
		})
	}
}

// Test ContainsRelativeTerm

func TestContainsRelativeTerm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		terms []string
		want  bool
	}{
		{
			name:  "Contains term",
			input: "2 days ago",
			terms: []string{"ago", "in"},
			want:  true,
		},
		{
			name:  "Contains at start",
			input: "ago was yesterday",
			terms: []string{"ago", "in"},
			want:  true,
		},
		{
			name:  "Contains at end",
			input: "2 days ago",
			terms: []string{"ago", "in"},
			want:  true,
		},
		{
			name:  "Case insensitive",
			input: "2 Days AGO",
			terms: []string{"ago", "in"},
			want:  true,
		},
		{
			name:  "No match",
			input: "tomorrow",
			terms: []string{"ago", "in"},
			want:  false,
		},
		{
			name:  "Empty input",
			input: "",
			terms: []string{"ago", "in"},
			want:  false,
		},
		{
			name:  "Empty terms",
			input: "ago",
			terms: []string{},
			want:  false,
		},
		{
			name:  "Empty term in list ignored",
			input: "test ago",
			terms: []string{"", "ago", ""},
			want:  true,
		},
		{
			name:  "Multiple terms one matches",
			input: "in 2 days",
			terms: []string{"ago", "in", "next"},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translations.ContainsRelativeTerm(tt.input, tt.terms)
			if got != tt.want {
				t.Errorf("ContainsRelativeTerm(%q, %v) = %v, want %v",
					tt.input, tt.terms, got, tt.want)
			}
		})
	}
}

// Test BuildTimeUnitPattern

func TestBuildTimeUnitPattern(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()

	tests := []struct {
		name          string
		languages     []*translations.Language
		shouldContain []string
		shouldNotBe   string
	}{
		{
			name:          "English units",
			languages:     []*translations.Language{english},
			shouldContain: []string{"day", "week", "month", "year"},
			shouldNotBe:   "",
		},
		{
			name:          "Spanish units",
			languages:     []*translations.Language{spanish},
			shouldContain: []string{"día", "semana", "mes", "año"},
			shouldNotBe:   "",
		},
		{
			name:          "Multiple languages",
			languages:     []*translations.Language{english, spanish},
			shouldContain: []string{"day", "día", "week", "semana"},
			shouldNotBe:   "",
		},
		{
			name:          "Empty languages",
			languages:     []*translations.Language{},
			shouldContain: []string{},
			shouldNotBe:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := translations.BuildTimeUnitPattern(tt.languages)

			if len(tt.languages) > 0 && pattern == "" {
				t.Error("BuildTimeUnitPattern() returned empty pattern")
			}

			for _, unit := range tt.shouldContain {
				if !strings.Contains(pattern, unit) {
					t.Errorf("BuildTimeUnitPattern() pattern should contain %q, got %q",
						unit, pattern)
				}
			}
		})
	}
}

// Test GetWeekdayPattern

func TestGetWeekdayPattern(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()

	tests := []struct {
		name          string
		languages     []*translations.Language
		shouldContain []string
		shouldNotBe   string
	}{
		{
			name:          "English weekdays",
			languages:     []*translations.Language{english},
			shouldContain: []string{"monday", "tuesday", "friday"},
			shouldNotBe:   "",
		},
		{
			name:          "Spanish weekdays",
			languages:     []*translations.Language{spanish},
			shouldContain: []string{"lunes", "martes", "viernes"},
			shouldNotBe:   "",
		},
		{
			name:          "Multiple languages",
			languages:     []*translations.Language{english, spanish},
			shouldContain: []string{"monday", "lunes", "friday", "viernes"},
			shouldNotBe:   "",
		},
		{
			name:          "Empty languages",
			languages:     []*translations.Language{},
			shouldContain: []string{},
			shouldNotBe:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := translations.GetWeekdayPattern(tt.languages)

			if len(tt.languages) > 0 && pattern == "" {
				t.Error("GetWeekdayPattern() returned empty pattern")
			}

			for _, weekday := range tt.shouldContain {
				if !strings.Contains(pattern, weekday) {
					t.Errorf("GetWeekdayPattern() pattern should contain %q, got %q",
						weekday, pattern)
				}
			}
		})
	}
}

// Test GetMonthPattern

func TestGetMonthPattern(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()

	tests := []struct {
		name          string
		languages     []*translations.Language
		shouldContain []string
		shouldNotBe   string
	}{
		{
			name:          "English months",
			languages:     []*translations.Language{english},
			shouldContain: []string{"january", "december", "jun"},
			shouldNotBe:   "",
		},
		{
			name:          "Spanish months",
			languages:     []*translations.Language{spanish},
			shouldContain: []string{"enero", "diciembre", "jun"},
			shouldNotBe:   "",
		},
		{
			name:          "Multiple languages",
			languages:     []*translations.Language{english, spanish},
			shouldContain: []string{"january", "enero", "december", "diciembre"},
			shouldNotBe:   "",
		},
		{
			name:          "Empty languages",
			languages:     []*translations.Language{},
			shouldContain: []string{},
			shouldNotBe:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := translations.GetMonthPattern(tt.languages)

			if len(tt.languages) > 0 && pattern == "" {
				t.Error("GetMonthPattern() returned empty pattern")
			}

			for _, month := range tt.shouldContain {
				if !strings.Contains(pattern, month) {
					t.Errorf("GetMonthPattern() pattern should contain %q, got %q",
						month, pattern)
				}
			}
		})
	}
}

// Test NormalizeTimeUnit

func TestNormalizeTimeUnit(t *testing.T) {
	english := translations.NewEnglishTranslation()
	spanish := translations.NewSpanishTranslation()
	french := translations.NewFrenchTranslation()

	tests := []struct {
		name      string
		input     string
		languages []*translations.Language
		want      string
	}{
		// English normalization
		{
			name:      "English day",
			input:     "day",
			languages: []*translations.Language{english},
			want:      "day",
		},
		{
			name:      "English days plural",
			input:     "days",
			languages: []*translations.Language{english},
			want:      "day",
		},
		{
			name:      "English week",
			input:     "week",
			languages: []*translations.Language{english},
			want:      "week",
		},
		{
			name:      "English month",
			input:     "month",
			languages: []*translations.Language{english},
			want:      "month",
		},
		{
			name:      "English year",
			input:     "year",
			languages: []*translations.Language{english},
			want:      "year",
		},
		{
			name:      "English fortnight",
			input:     "fortnight",
			languages: []*translations.Language{english},
			want:      "fortnight",
		},
		{
			name:      "English decade",
			input:     "decade",
			languages: []*translations.Language{english},
			want:      "decade",
		},
		{
			name:      "English hour",
			input:     "hour",
			languages: []*translations.Language{english},
			want:      "hour",
		},
		{
			name:      "English minute",
			input:     "minute",
			languages: []*translations.Language{english},
			want:      "minute",
		},
		{
			name:      "English second",
			input:     "second",
			languages: []*translations.Language{english},
			want:      "second",
		},

		// Spanish normalization
		{
			name:      "Spanish día",
			input:     "día",
			languages: []*translations.Language{spanish},
			want:      "day",
		},
		{
			name:      "Spanish días",
			input:     "días",
			languages: []*translations.Language{spanish},
			want:      "day",
		},
		{
			name:      "Spanish semana",
			input:     "semana",
			languages: []*translations.Language{spanish},
			want:      "week",
		},
		{
			name:      "Spanish mes",
			input:     "mes",
			languages: []*translations.Language{spanish},
			want:      "month",
		},
		{
			name:      "Spanish año",
			input:     "año",
			languages: []*translations.Language{spanish},
			want:      "year",
		},

		// French normalization
		{
			name:      "French jour",
			input:     "jour",
			languages: []*translations.Language{french},
			want:      "day",
		},
		{
			name:      "French jours",
			input:     "jours",
			languages: []*translations.Language{french},
			want:      "day",
		},

		// Multi-language
		{
			name:      "Multi-language Spanish día",
			input:     "día",
			languages: []*translations.Language{english, spanish},
			want:      "day",
		},

		// Edge cases
		{
			name:      "Case insensitive",
			input:     "DAY",
			languages: []*translations.Language{english},
			want:      "day",
		},
		{
			name:      "With whitespace",
			input:     "  day  ",
			languages: []*translations.Language{english},
			want:      "day",
		},
		{
			name:      "Unknown unit",
			input:     "unknown",
			languages: []*translations.Language{english},
			want:      "unknown",
		},
		{
			name:      "Empty input",
			input:     "",
			languages: []*translations.Language{english},
			want:      "",
		},
		{
			name:      "No languages",
			input:     "day",
			languages: []*translations.Language{},
			want:      "day",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translations.NormalizeTimeUnit(tt.input, tt.languages)
			if got != tt.want {
				t.Errorf("NormalizeTimeUnit(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
