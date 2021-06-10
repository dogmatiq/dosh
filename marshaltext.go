package dosh

import (
	"errors"
	"fmt"
	"strings"

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
	str := strings.TrimSpace(string(text))
	parts := strings.SplitN(str, " ", 2)

	if len(parts) != 2 {
		return errors.New("cannot unmarshal amount from text representation: data must have currency and magnitude components")
	}

	m, err := decimal.NewFromString(parts[1])
	if err != nil {
		return fmt.Errorf("cannot unmarshal amount from text representation: %w", err)
	}

	a.cur = strings.ToUpper(parts[0])
	a.mag = m

	return nil
}
