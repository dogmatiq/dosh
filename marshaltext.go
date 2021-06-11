package dosh

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/dogmatiq/dosh/internal/currency"
	"github.com/shopspring/decimal"
)

// MarshalText mashals an amount to its text representation.
func (a Amount) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

// UnmarshalText unmarshals an amount from its text representation.
//
// NOTE: In order to comply with Go's encoding.TextUnmarshaler interface, this
// method mutates the internals of a, violating Amount's immutability guarantee.
func (a *Amount) UnmarshalText(text []byte) error {
	n := bytes.IndexRune(text, ' ')
	if n == -1 {
		return errors.New("cannot unmarshal amount from text representation: data must have currency and magnitude components separated by a single space")
	}

	c := string(text[:n])
	m := string(text[n+1:])

	if err := currency.ValidateCode(c); err != nil {
		return fmt.Errorf("cannot unmarshal amount from text representation: %w", err)
	}

	d, err := decimal.NewFromString(m)
	if err != nil {
		if strings.TrimSpace(m) == "" {
			// Provide a slightly less-cryptic error message when the magnitude
			// consists solely of whitespace. Otherwise, we get an error like
			// "can't convert  to decimal".
			return errors.New("cannot unmarshal amount from text representation: cannot parse magnitude")
		}

		return fmt.Errorf("cannot unmarshal amount from text representation: %w", err)
	}

	a.cur = c
	a.mag = d

	return nil
}
