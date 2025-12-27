package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSafeRandomString generates a cryptographically secure random string
// of the specified length using base64 URL-safe encoding.
// The actual output length may be slightly longer due to base64 encoding.
func GenerateSafeRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive, got: %d", length)
	}

	// Calculate the number of bytes needed
	// base64 encoding produces 4 characters for every 3 bytes
	numBytes := (length * 3) / 4
	if numBytes == 0 {
		numBytes = 1
	}

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to URL-safe base64 and trim to desired length
	encoded := base64.URLEncoding.EncodeToString(randomBytes)
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
