package validation

import (
	"fmt"
	"strconv"
)

func ValidateGaugeValue(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func ValidateCounterValue(value string) error {
	_, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}
