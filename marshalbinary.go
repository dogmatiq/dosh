package dosh

import (
	"bytes"
	"errors"
	"fmt"

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
		[]byte{byte(n)},
		[]byte(c)...,
	)

	m, err := a.mag.MarshalBinary()
	if err != nil {
		// CODE COVERAGE: It does not appear this branch can currently be
		// reached as decimal.Decimal never returns a non-nil error.
		return nil, fmt.Errorf("cannot marshal amount to binary representation: %w", err)
	}

	data = append(data, m...)

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

	n := int(data[0])
	data = data[1:]

	if n == 0 {
		return errors.New("cannot unmarshal amount from binary representation: currency component is empty")
	}

	if len(data) < n+4 {
		// Note: +4 is a workaround to https://github.com/shopspring/decimal/issues/231.
		return errors.New("cannot unmarshal amount from binary representation: data is shorter than expected")
	}

	c := data[:n]
	data = data[n:]

	var m decimal.Decimal
	if err := m.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("cannot unmarshal amount from binary representation: %w", err)
	}

	a.cur = bytes.ToUpper(c)
	a.mag = m

	return nil
}
