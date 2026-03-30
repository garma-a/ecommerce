package products

import (
	"strings"
	"unicode/utf8"
)

const (
	minNameLen = 2
	maxNameLen = 120
	minPrice   = 0.01
	maxPrice   = 1_000_000
)

// ValidateProductInput validates and normalizes product input
func ValidateProductInput(name string, price float64) (string, error) {
	// Normalize name: trim whitespace
	trimmed := strings.TrimSpace(name)

	// Validate name length
	if trimmed == "" {
		return "", NewValidationError("name", "product name cannot be empty")
	}

	nameLen := utf8.RuneCountInString(trimmed)
	if nameLen < minNameLen {
		return "", NewValidationError("name", "product name must be at least 2 characters")
	}
	if nameLen > maxNameLen {
		return "", NewValidationError("name", "product name cannot exceed 120 characters")
	}

	// Validate price
	if price < minPrice {
		return "", NewValidationError("price", "product price must be at least 0.01")
	}
	if price > maxPrice {
		return "", NewValidationError("price", "product price cannot exceed 1,000,000")
	}

	return trimmed, nil
}
