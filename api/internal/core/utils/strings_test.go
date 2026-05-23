package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAllIndicesOfChar tests finding all indices of a character in a string
func TestAllIndicesOfChar(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		charToFind string
		expected   []int
	}{
		{
			name:       "multiple occurrences",
			str:        "/path/to/file",
			charToFind: "/",
			expected:   []int{0, 5, 8},
		},
		{
			name:       "single occurrence",
			str:        "hello world",
			charToFind: "w",
			expected:   []int{6},
		},
		{
			name:       "no occurrences",
			str:        "hello world",
			charToFind: "z",
			expected:   nil,
		},
		{
			name:       "empty string",
			str:        "",
			charToFind: "a",
			expected:   nil,
		},
		{
			name:       "all same character",
			str:        "aaaa",
			charToFind: "a",
			expected:   []int{0, 1, 2, 3},
		},
		{
			name:       "special characters",
			str:        "foo.bar.baz",
			charToFind: ".",
			expected:   []int{3, 7},
		},
		{
			name:       "unicode characters",
			str:        "hello 世界 world",
			charToFind: " ",
			expected:   []int{5, 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AllIndicesOfChar(tt.str, tt.charToFind)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetEnv tests environment variable retrieval with fallback
func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		fallback     string
		envValue     string
		setEnv       bool
		expectPanic  bool
		expectedValue string
	}{
		{
			name:          "env var exists",
			key:           "TEST_VAR_EXISTS",
			fallback:      "default",
			envValue:      "actual_value",
			setEnv:        true,
			expectPanic:   false,
			expectedValue: "actual_value",
		},
		{
			name:          "env var missing with fallback",
			key:           "TEST_VAR_MISSING",
			fallback:      "fallback_value",
			envValue:      "",
			setEnv:        false,
			expectPanic:   false,
			expectedValue: "fallback_value",
		},
		{
			name:        "env var missing without fallback",
			key:         "TEST_VAR_MISSING_NO_FALLBACK",
			fallback:    "",
			envValue:    "",
			setEnv:      false,
			expectPanic: true,
		},
		{
			name:          "empty env var with fallback",
			key:           "TEST_VAR_EMPTY",
			fallback:      "fallback_value",
			envValue:      "",
			setEnv:        true,
			expectPanic:   false,
			expectedValue: "fallback_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before and after
			defer os.Unsetenv(tt.key)
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					GetEnv(tt.key, tt.fallback)
				})
			} else {
				result := GetEnv(tt.key, tt.fallback)
				assert.Equal(t, tt.expectedValue, result)
			}
		})
	}
}

// TestExists tests checking if a value exists in a slice
func TestExists(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		val      string
		expected bool
	}{
		{
			name:     "value exists",
			slice:    []string{"apple", "banana", "cherry"},
			val:      "banana",
			expected: true,
		},
		{
			name:     "value does not exist",
			slice:    []string{"apple", "banana", "cherry"},
			val:      "grape",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			val:      "apple",
			expected: false,
		},
		{
			name:     "empty value in slice",
			slice:    []string{"", "apple", "banana"},
			val:      "",
			expected: true,
		},
		{
			name:     "case sensitive match",
			slice:    []string{"Apple", "Banana", "Cherry"},
			val:      "apple",
			expected: false,
		},
		{
			name:     "exact match",
			slice:    []string{"apple", "banana", "cherry"},
			val:      "apple",
			expected: true,
		},
		{
			name:     "value at end",
			slice:    []string{"apple", "banana", "cherry"},
			val:      "cherry",
			expected: true,
		},
		{
			name:     "duplicate values",
			slice:    []string{"apple", "banana", "apple"},
			val:      "apple",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Exists(tt.slice, tt.val)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEscapeInvalidCharacters tests escaping single quotes in strings
func TestEscapeInvalidCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single quote",
			input:    "O'Brien",
			expected: "O\\'Brien",
		},
		{
			name:     "multiple single quotes",
			input:    "It's a beautiful day's work",
			expected: "It\\'s a beautiful day\\'s work",
		},
		{
			name:     "no single quotes",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only single quote",
			input:    "'",
			expected: "\\'",
		},
		{
			name:     "SQL injection attempt",
			input:    "'; DROP TABLE users; --",
			expected: "\\'; DROP TABLE users; --",
		},
		{
			name:     "string with backslashes",
			input:    "path\\to\\file",
			expected: "path\\to\\file",
		},
		{
			name:     "unicode with single quote",
			input:    "世界's best",
			expected: "世界\\'s best",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeInvalidCharacters(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestUnescapeInvalidCharacters tests unescaping single quotes in strings
func TestUnescapeInvalidCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped single quote",
			input:    "O\\'Brien",
			expected: "O'Brien",
		},
		{
			name:     "multiple escaped single quotes",
			input:    "It\\'s a beautiful day\\'s work",
			expected: "It's a beautiful day's work",
		},
		{
			name:     "no escaped quotes",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only escaped single quote",
			input:    "\\'",
			expected: "'",
		},
		{
			name:     "directory with apostrophe",
			input:    "Krazy House Jan \\'07",
			expected: "Krazy House Jan '07",
		},
		{
			name:     "unicode with escaped quote",
			input:    "世界\\'s best",
			expected: "世界's best",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnescapeInvalidCharacters(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
