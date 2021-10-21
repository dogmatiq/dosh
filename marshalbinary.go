package dosh

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/dosh/internal/currency"
	"github.com/shopspring/decimal"
)

// MarshalBinary mashals an amount to its binary representation.
func (a Amount) MarshalBinary() ([]byte, error) {
	c := a.CurrencyCode()
	n := len(c)

	if n > 255 {
		return nil, fmt.Errorf("cannot marshal amount to binary representation: currency code is %d characters long, maximum is 255", n)
	}

	data := append(
		[]byte{byte(n)}, // 1-byte length of currency code
		[]byte(c)...,    // currency code itself
	)

	m, err := a.mag.MarshalBinary()
	if err != nil {
		// CODE COVERAGE: It does not appear this branch can currently be
		// reached as decimal.Decimal.MarshalBinary() never returns a non-nil
		// error.
		return nil, fmt.Errorf("cannot marshal amount to binary representation: %w", err)
	}

	data = append(data, m...) // GOB encoded decimal

	return data, nil
}

// UnmarshalBinary unmarshals an amount from its binary representation.
//
// NOTE: In order to comply with Go's encoding.BinaryUnmarshaler interface, this
// method mutates the internals of a, violating Amount's immutability guarantee.
func (a *Amount) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("cannot unmarshal amount from binary representation: data is empty")
	}

	// Read 1-byte length of currency code.
	n := int(data[0])
	data = data[1:]

	// Unmarshal the currency code.
	c := string(data[:n])
	data = data[n:]

	if err := currency.ValidateCode(c); err != nil {
		return fmt.Errorf("cannot unmarshal amount from binary representation: %w", err)
	}

	// Unmarshal the magnitude component.
	var m decimal.Decimal
	if err := m.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("cannot unmarshal amount from binary representation: %w", err)
	}

	a.cur = c
	a.mag = m

	return nil
}
