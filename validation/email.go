package validation

import (
	"fmt"
	"strings"
)

// ValidateEmail performs basic email validation
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email validation
	if len(email) < 3 || !strings.Contains(email, "@") {
		return fmt.Errorf("please enter a valid email address")
	}

	return nil
}
