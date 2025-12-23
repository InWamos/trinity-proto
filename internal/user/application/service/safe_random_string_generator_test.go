package service_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/InWamos/trinity-proto/internal/user/application/service"
)

func TestGenerateSafeRandomString(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{
			name:        "Valid length 32",
			length:      32,
			expectError: false,
		},
		{
			name:        "Valid length 16",
			length:      16,
			expectError: false,
		},
		{
			name:        "Valid length 64",
			length:      64,
			expectError: false,
		},
		{
			name:        "Minimum length 1",
			length:      1,
			expectError: false,
		},
		{
			name:        "Zero length",
			length:      0,
			expectError: true,
		},
		{
			name:        "Negative length",
			length:      -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.GenerateSafeRandomString(tt.length)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tt.length {
				t.Errorf("expected length %d, got %d", tt.length, len(result))
			}

			// Verify it's valid base64 URL-safe encoding
			if !isURLSafeBase64(result) {
				t.Errorf("result contains non-URL-safe base64 characters: %s", result)
			}
		})
	}
}

func TestGenerateSafeRandomString_Uniqueness(t *testing.T) {
	const iterations = 100
	const length = 32

	seen := make(map[string]bool)

	for i := range iterations {
		result, err := service.GenerateSafeRandomString(length)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}

		if seen[result] {
			t.Errorf("duplicate random string generated: %s", result)
		}
		seen[result] = true
	}

	if len(seen) != iterations {
		t.Errorf("expected %d unique strings, got %d", iterations, len(seen))
	}
}

func TestGenerateSafeRandomString_NoSpecialCharacters(t *testing.T) {
	result, err := service.GenerateSafeRandomString(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// URL-safe base64 should not contain + or /
	if strings.Contains(result, "+") {
		t.Errorf("result contains '+' character, not URL-safe: %s", result)
	}
	if strings.Contains(result, "/") {
		t.Errorf("result contains '/' character, not URL-safe: %s", result)
	}
}

func TestGenerateSafeRandomString_ConsistentLength(t *testing.T) {
	lengths := []int{8, 16, 24, 32, 48, 64, 128}

	for _, length := range lengths {
		t.Run(string(rune(length)), func(t *testing.T) {
			for i := range 10 {
				result, err := service.GenerateSafeRandomString(length)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if len(result) != length {
					t.Errorf("iteration %d: expected length %d, got %d", i, length, len(result))
				}
			}
		})
	}
}

// isURLSafeBase64 checks if a string contains only URL-safe base64 characters.
func isURLSafeBase64(s string) bool {
	// URL-safe base64 alphabet: A-Z, a-z, 0-9, -, _
	for _, char := range s {
		if (char < 'A' || char > 'Z') &&
			(char < 'a' || char > 'z') &&
			(char < '0' || char > '9') &&
			char != '-' && char != '_' {
			return false
		}
	}
	return true
}

func BenchmarkGenerateSafeRandomString(b *testing.B) {
	lengths := []int{16, 32, 64, 128}

	for _, length := range lengths {
		b.Run(string(rune(length)), func(b *testing.B) {
			for range b.N {
				_, err := service.GenerateSafeRandomString(length)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestGenerateSafeRandomString_Base64Decoding(t *testing.T) {
	result, err := service.GenerateSafeRandomString(32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Pad the string if necessary for base64 decoding
	padded := result
	if mod := len(padded) % 4; mod != 0 {
		padded += strings.Repeat("=", 4-mod)
	}

	// Should be decodable as base64
	_, err = base64.URLEncoding.DecodeString(padded)
	if err != nil {
		t.Errorf("generated string is not valid base64: %v", err)
	}
}
