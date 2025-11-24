package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatNumber formats number with thousand separators (using dot)
func FormatNumber(num interface{}, decimals int) string {
	var numFloat float64
	switch v := num.(type) {
	case float64:
		numFloat = v
	case int:
		numFloat = float64(v)
	case int64:
		numFloat = float64(v)
	case string:
		var err error
		numFloat, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return "0"
		}
	default:
		return "0"
	}

	// Format with decimals
	format := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(format, numFloat)

	// Split integer and decimal parts
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Add thousand separators (dots) to integer part from right to left
	// Reverse the string, add dots every 3 digits, then reverse back
	reversed := reverseString(integerPart)
	var formattedReversed strings.Builder
	for i, char := range reversed {
		if i > 0 && i%3 == 0 {
			formattedReversed.WriteString(".")
		}
		formattedReversed.WriteRune(char)
	}
	integerPart = reverseString(formattedReversed.String())

	// Build final result
	var result strings.Builder
	result.WriteString(integerPart)

	// Add decimal part if exists
	if decimalPart != "" {
		result.WriteString(",")
		result.WriteString(decimalPart)
	}

	return result.String()
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

