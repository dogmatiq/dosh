package currency

import (
	"errors"
	"fmt"
)

// ValidateCode returns an error if s is not a valid currency code.
//
// A valid currency code is a minimum of 3 characters long, and consists
// entirely of uppercase ASCII letters.
func ValidateCode(c string) error {
	if len(c) == 0 {
		return errors.New("currency code is empty, codes must consist only of 3 or more uppercase ASCII letters")
	}

	if len(c) >= 3 && isUppercaseASCII(c) {
		return nil
	}

	return fmt.Errorf(
		"currency code (%s) is invalid, codes must consist only of 3 or more uppercase ASCII letters",
		c,
	)
}

// isUppercaseASCII returns true if c consists only of uppercase ASCII letters.
func isUppercaseASCII(c string) bool {
	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}
