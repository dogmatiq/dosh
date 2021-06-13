package protomoney

import (
	"fmt"

	"github.com/dogmatiq/dosh/internal/currency"
	"google.golang.org/genproto/googleapis/type/money"
)

// nanosPerUnit is the number of "nano units" in each unit.
const nanosPerUnit = 1_000_000_000

// Validate returns an error if m is invalid.
func Validate(m *money.Money) error {
	if err := currency.ValidateCode(m.CurrencyCode); err != nil {
		return err
	}

	return checkSignsAgree(m)
}

// Normalize validates m, and returns its normalized version if it is
// valid.
//
// Within a normalized value the nanos component is guaranteed to be less than
// one whole unit.
//
// m itself is never mutated. If it is already normalized it is returned
// unchanged; otherwise a normalized clone is returned.
func Normalize(m *money.Money) (*money.Money, error) {
	if err := Validate(m); err != nil {
		return nil, err
	}

	return normalize(m), nil
}

// isNormalized returns true if the units and nanos components of m are already
// normalized.
func isNormalized(m *money.Money) bool {
	assertSignsAgree(m)
	return -nanosPerUnit < m.Nanos && m.Nanos < nanosPerUnit
}

// normalizeInPlace normalizes the units and nanos components of m such that
// m.Nanos is guaranteed to contain less than one whole unit.
func normalizeInPlace(m *money.Money) {
	if !isNormalized(m) {
		m.Units += int64(m.Nanos) / nanosPerUnit
		m.Nanos %= nanosPerUnit
	}
}

// normalize returns m with normalized units and nanos components, or panics
// if unable to do so.
func normalize(m *money.Money) *money.Money {
	if isNormalized(m) {
		return m
	}

	return &money.Money{
		CurrencyCode: m.CurrencyCode,
		Units:        m.Units + (int64(m.Nanos) / nanosPerUnit),
		Nanos:        m.Nanos % nanosPerUnit,
	}
}

// assertSameCurrency panics if a and b do not have the same currency.
func assertSameCurrency(a, b *money.Money) {
	if a.CurrencyCode != b.CurrencyCode {
		panic(fmt.Sprintf(
			"can not operate on amounts in differing currencies (%s vs %s)",
			a.CurrencyCode,
			b.CurrencyCode,
		))
	}
}

// assertSignsAgree panics if the signs of m.Units and m.Nanos do not agree.
func assertSignsAgree(m *money.Money) {
	if err := checkSignsAgree(m); err != nil {
		panic(err)
	}
}

// checkSignsAgree returns an error if the signs of m.Units and m.Nanos do not
// agree.
func checkSignsAgree(m *money.Money) error {
	if (m.Units > 0 && m.Nanos < 0) || (m.Units < 0 && m.Nanos > 0) {
		return fmt.Errorf(
			"sign of units component (%d) does not agree with sign of nanos component (%d)",
			m.Units,
			m.Nanos,
		)
	}

	return nil
}
