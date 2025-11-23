package helpers

import (
	"fmt"
	"strconv"
)

// FormatNumber formats number with thousand separators
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

	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, numFloat)
}

