package protomoney

import (
	"google.golang.org/genproto/googleapis/type/money"
)

// IsZero returns true if m has a magnitude of zero.
func IsZero(m *money.Money) bool {
	return m.Units == 0 && m.Nanos == 0
}

// IsPositive returns true m has a positive magnitude.
func IsPositive(m *money.Money) bool {
	assertSignsAgree(m)
	return m.Units > 0 || m.Nanos > 0
}

// IsNegative returns true m has a negative magnitude.
func IsNegative(m *money.Money) bool {
	assertSignsAgree(m)
	return m.Units < 0 || m.Nanos < 0
}

// Cmp compares a to b and returns a C-style comparison result.
//
// It panics if a and b do not use the same currency.
//
// If a < b then c is negative.
// If a > b then c is positive.
// Otherwise; a == b and c is zero.
func Cmp(a, b *money.Money) (c int) {
	assertSameCurrency(a, b)

	a = normalize(a)
	b = normalize(b)

	if a.Units < b.Units {
		return -1
	}

	if a.Units > b.Units {
		return +1
	}

	return int(a.Nanos - b.Nanos)
}

// Equal returns true if a and b have the same magnitude.
//
// It panics if a and b do not use the same currency.
//
// To check equality between two amounts that may have differing currencies, use
// IdenticalTo() instead.
func EqualTo(a, b *money.Money) bool {
	return Cmp(a, b) == 0
}

// Identical returns true if a and b use the same currency and have the same
// magnitude.
//
// For general comparisons that are expected to be in the same currency, use
// EqualTo() instead.
func IdenticalTo(a, b *money.Money) bool {
	if a.CurrencyCode != b.CurrencyCode {
		return false
	}

	return Cmp(a, b) == 0
}

// LessThan returns true if a < b.
//
// It panics if a and b do not use the same currency.
func LessThan(a, b *money.Money) bool {
	return Cmp(a, b) < 0
}

// LessThanOrEqualTo returns true if a <= b.
//
// It panics if a and b do not use the same currency.
func LessThanOrEqualTo(a, b *money.Money) bool {
	return Cmp(a, b) <= 0
}

// GreaterThan returns true if a > b.
//
// It panics if a and b do not use the same currency.
func GreaterThan(a, b *money.Money) bool {
	return Cmp(a, b) > 0
}

// GreaterThanOrEqualTo returns true if a >= b.
//
// It panics if a and b do not use the same currency.
func GreaterThanOrEqualTo(a, b *money.Money) bool {
	return Cmp(a, b) >= 0
}

// LexicallyLessThan returns true if a should appear before b in a sorted list.
//
// There is no requirement that a and b use the same currency.
func LexicallyLessThan(a, b *money.Money) bool {
	if a.CurrencyCode == b.CurrencyCode {
		return LessThan(a, b)
	}

	return a.CurrencyCode < b.CurrencyCode
}

// Min returns the smallest of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Min(amounts ...*money.Money) *money.Money {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		if LessThan(b, a) {
			a = b
		}
	}

	return a
}

// Max returns the largest of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Max(amounts ...*money.Money) *money.Money {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		if GreaterThan(b, a) {
			a = b
		}
	}

	return a
}
