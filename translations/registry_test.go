package translations_test

import (
	"testing"

	"github.com/coredds/godateparser/translations"
)

// Test Registry Creation and Registration

func TestNewRegistry(t *testing.T) {
	registry := translations.NewRegistry()

	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	// Should have English registered by default
	eng := registry.Get("en")
	if eng == nil {
		t.Error("NewRegistry() should register English by default")
	}

	if eng.Code != "en" {
		t.Errorf("Default language code = %q, want 'en'", eng.Code)
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := translations.NewRegistry()

	// Register Spanish
	spanish := translations.NewSpanishTranslation()
	registry.Register(spanish)

	// Retrieve and verify
	retrieved := registry.Get("es")
	if retrieved == nil {
		t.Fatal("Register() failed to register Spanish")
	}

	if retrieved.Code != "es" {
		t.Errorf("Retrieved language code = %q, want 'es'", retrieved.Code)
	}

	if retrieved.Name != "Spanish" {
		t.Errorf("Retrieved language name = %q, want 'Spanish'", retrieved.Name)
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := translations.NewRegistry()

	tests := []struct {
		name     string
		code     string
		wantCode string
		wantName string
	}{
		{
			name:     "English (default)",
			code:     "en",
			wantCode: "en",
			wantName: "English",
		},
		{
			name:     "Unknown language returns default",
			code:     "xx",
			wantCode: "en",
			wantName: "English",
		},
		{
			name:     "Empty code returns default",
			code:     "",
			wantCode: "en",
			wantName: "English",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := registry.Get(tt.code)
			if lang == nil {
				t.Fatal("Get() returned nil")
			}

			if lang.Code != tt.wantCode {
				t.Errorf("Get(%q) code = %q, want %q", tt.code, lang.Code, tt.wantCode)
			}

			if lang.Name != tt.wantName {
				t.Errorf("Get(%q) name = %q, want %q", tt.code, lang.Name, tt.wantName)
			}
		})
	}
}

func TestRegistry_GetMultiple(t *testing.T) {
	registry := translations.NewRegistry()
	registry.Register(translations.NewSpanishTranslation())
	registry.Register(translations.NewFrenchTranslation())

	tests := []struct {
		name       string
		codes      []string
		wantCodes  []string
		wantLength int
	}{
		{
			name:       "Single language",
			codes:      []string{"en"},
			wantCodes:  []string{"en"},
			wantLength: 1,
		},
		{
			name:       "Multiple languages",
			codes:      []string{"es", "fr", "en"},
			wantCodes:  []string{"es", "fr", "en"},
			wantLength: 3,
		},
		{
			name:       "Mixed valid and invalid",
			codes:      []string{"es", "xx", "en"},
			wantCodes:  []string{"es", "en"},
			wantLength: 2,
		},
		{
			name:       "Empty codes returns default",
			codes:      []string{},
			wantCodes:  []string{"en"},
			wantLength: 1,
		},
		{
			name:       "All invalid returns default",
			codes:      []string{"xx", "yy", "zz"},
			wantCodes:  []string{"en"},
			wantLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			langs := registry.GetMultiple(tt.codes)

			if len(langs) != tt.wantLength {
				t.Errorf("GetMultiple() length = %d, want %d", len(langs), tt.wantLength)
			}

			for i, want := range tt.wantCodes {
				if i >= len(langs) {
					t.Errorf("GetMultiple() missing language at index %d", i)
					continue
				}
				if langs[i].Code != want {
					t.Errorf("GetMultiple() code[%d] = %q, want %q", i, langs[i].Code, want)
				}
			}
		})
	}
}

func TestRegistry_DetectLanguage(t *testing.T) {
	registry := translations.NewRegistry()
	registry.Register(translations.NewSpanishTranslation())
	registry.Register(translations.NewFrenchTranslation())
	registry.Register(translations.NewGermanTranslation())
	registry.Register(translations.NewPortugueseTranslation())
	registry.Register(translations.NewItalianTranslation())
	registry.Register(translations.NewDutchTranslation())
	registry.Register(translations.NewRussianTranslation())
	registry.Register(translations.NewChineseTranslation())
	registry.Register(translations.NewJapaneseTranslation())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		// English detection
		{
			name:  "English month",
			input: "next Thursday",
			want:  "en",
		},
		{
			name:  "English weekday",
			input: "next Wednesday",
			want:  "en",
		},

		// Spanish detection
		{
			name:  "Spanish month",
			input: "31 diciembre 2024",
			want:  "es",
		},
		{
			name:  "Spanish weekday",
			input: "próximo lunes",
			want:  "es",
		},
		{
			name:  "Spanish relative",
			input: "ayer",
			want:  "es",
		},

		// French detection
		{
			name:  "French month",
			input: "31 décembre 2024",
			want:  "fr",
		},
		{
			name:  "French weekday",
			input: "prochain lundi",
			want:  "fr",
		},
		{
			name:  "French relative",
			input: "hier",
			want:  "fr",
		},

		// German detection
		{
			name:  "German month",
			input: "31 Dezember 2024",
			want:  "de",
		},
		{
			name:  "German weekday",
			input: "nächsten Montag",
			want:  "de",
		},

		// Portuguese detection
		{
			name:  "Portuguese month",
			input: "31 dezembro 2024",
			want:  "pt",
		},
		{
			name:  "Portuguese relative",
			input: "ontem",
			want:  "pt",
		},

		// Italian detection
		{
			name:  "Italian month",
			input: "31 dicembre 2024",
			want:  "it",
		},
		{
			name:  "Italian relative",
			input: "ieri",
			want:  "it",
		},

		// Dutch detection
		{
			name:  "Dutch month",
			input: "31 mei 2024",
			want:  "nl",
		},
		{
			name:  "Dutch relative",
			input: "vandaag",
			want:  "nl",
		},

		// Russian detection
		{
			name:  "Russian month",
			input: "31 декабря 2024",
			want:  "ru",
		},
		{
			name:  "Russian weekday",
			input: "понедельник",
			want:  "ru",
		},
		{
			name:  "Russian relative",
			input: "вчера",
			want:  "ru",
		},

		// Chinese detection
		{
			name:  "Chinese weekday",
			input: "星期一",
			want:  "zh",
		},
		{
			name:  "Chinese relative",
			input: "昨天",
			want:  "zh",
		},
		{
			name:  "Chinese week term",
			input: "下周",
			want:  "zh",
		},

		// Japanese detection
		{
			name:  "Japanese month",
			input: "来月",
			want:  "ja",
		},
		{
			name:  "Japanese weekday",
			input: "月曜日",
			want:  "ja",
		},
		{
			name:  "Japanese relative",
			input: "先週",
			want:  "ja",
		},

		// Default fallback (Note: may vary due to scoring, these test basic fallback behavior)
		{
			name:  "Empty input returns default",
			input: "",
			want:  "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := registry.DetectLanguage(tt.input)
			if got != tt.want {
				t.Errorf("DetectLanguage(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRegistry_SupportedLanguages(t *testing.T) {
	registry := translations.NewRegistry()

	// Initial state - should have at least English
	langs := registry.SupportedLanguages()
	if len(langs) == 0 {
		t.Fatal("SupportedLanguages() returned empty slice")
	}

	hasEnglish := false
	for _, code := range langs {
		if code == "en" {
			hasEnglish = true
			break
		}
	}
	if !hasEnglish {
		t.Error("SupportedLanguages() should include 'en'")
	}

	// Add more languages
	registry.Register(translations.NewSpanishTranslation())
	registry.Register(translations.NewFrenchTranslation())

	langs = registry.SupportedLanguages()
	if len(langs) < 3 {
		t.Errorf("SupportedLanguages() length = %d, want >= 3", len(langs))
	}

	// Check that Spanish and French are included
	expected := map[string]bool{"en": true, "es": true, "fr": true}
	for _, code := range langs {
		if _, ok := expected[code]; ok {
			delete(expected, code)
		}
	}

	if len(expected) > 0 {
		t.Errorf("SupportedLanguages() missing languages: %v", expected)
	}
}

func TestGlobalRegistry(t *testing.T) {
	// Test that the global registry is initialized
	lang := translations.GetLanguage("en")
	if lang == nil {
		t.Fatal("GetLanguage('en') returned nil")
	}

	if lang.Code != "en" {
		t.Errorf("GetLanguage('en') code = %q, want 'en'", lang.Code)
	}

	// Test language detection
	detected := translations.DetectLanguage("December 31")
	if detected == "" {
		t.Error("DetectLanguage() returned empty string")
	}

	// Test supported languages
	supported := translations.SupportedLanguages()
	if len(supported) == 0 {
		t.Fatal("SupportedLanguages() returned empty slice")
	}

	// Should have all 10 languages
	expectedLangs := []string{"en", "es", "pt", "fr", "de", "it", "nl", "ru", "zh", "ja"}
	found := make(map[string]bool)
	for _, code := range supported {
		found[code] = true
	}

	for _, expected := range expectedLangs {
		if !found[expected] {
			t.Errorf("SupportedLanguages() missing %q", expected)
		}
	}
}
