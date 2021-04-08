package dosh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// MarshalBinary mashals an amount to its binary representation.
func (a Amount) MarshalBinary() ([]byte, error) {
	c := a.CurrencyCode()
	n := len(c)

	if n > 255 {
		return nil, errors.New("cannot marshal amount to binary representation: currency code is too long")
	}

	data := append(
		[]byte{byte(n)},
		[]byte(c)...,
	)

	m, err := a.magnitude.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal amount to binary representation: %w", err)
	}

	data = append(data, m...)

	return data, nil
}

// UnmarshalBinary unmarshals an amount from its protocol buffers
// representation.
func (a *Amount) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("cannot unmarshal amount from binary representation: data is empty")
	}

	n := int(data[0])
	data = data[1:]

	if n == 0 {
		return errors.New("cannot unmarshal amount from binary representation: currency is empty")
	}

	if len(data) < n {
		return errors.New("cannot unmarshal amount from binary representation: data is shorted than expected")
	}

	c := string(data[:n])
	data = data[n:]

	var m decimal.Decimal
	if err := m.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("cannot unmarshal amount from binary representation: %w", err)
	}

	a.currency = strings.ToUpper(c)
	a.magnitude = m

	return nil
}
