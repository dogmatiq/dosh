package dosh

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

// MarshalText mashals an amount to its text representation.
func (a Amount) MarshalText() (text []byte, err error) {
	return []byte(a.String()), nil
}

// UnmarshalText unmarshals an amount from its protocol buffers representation.
func (a *Amount) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	parts := strings.SplitN(str, " ", 2)

	if len(parts) != 2 {
		return errors.New("cannot unmarshal amount from text representation: amount must have magnitude and currency components")
	}

	if parts[0] == "" {
		return errors.New("cannot unmarshal amount from text representation: currency code must not be empty")
	}

	m, err := decimal.NewFromString(parts[1])
	if err != nil {
		return err
	}

	a.currency = strings.ToUpper(parts[0])
	a.magnitude = m

	return nil
}
