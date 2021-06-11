package protomoney

import (
	"google.golang.org/genproto/googleapis/type/money"
)

// Abs returns the absolute value of this amount.
//
// That is, if m is negative, it returns its inverse (a positive magnitude),
// otherwise it returns m unchanged.
func Abs(m *money.Money) *money.Money {
	if IsPositive(m) {
		return m
	}

	return &money.Money{
		CurrencyCode: m.CurrencyCode,
		Units:        -m.Units,
		Nanos:        -m.Nanos,
	}
}

// Neg returns -m.
func Neg(m *money.Money) *money.Money {
	return &money.Money{
		CurrencyCode: m.CurrencyCode,
		Units:        -m.Units,
		Nanos:        -m.Nanos,
	}
}

// Add returns a + b.
//
// It panics if a and b do not use the same currency.
func Add(a, b *money.Money) *money.Money {
	assertSameCurrency(a, b)

	m := &money.Money{
		CurrencyCode: a.CurrencyCode,
		Units:        a.Units + b.Units,
		Nanos:        a.Nanos + b.Nanos,
	}

	normalizeInPlace(m)

	return m
}

// Sub returns a - b.
//
// It panics if a and b do not use the same currency.
func Sub(a, b *money.Money) *money.Money {
	assertSameCurrency(a, b)

	m := &money.Money{
		CurrencyCode: a.CurrencyCode,
		Units:        a.Units - b.Units,
		Nanos:        a.Nanos - b.Nanos,
	}

	normalizeInPlace(m)

	return m
}

// Sum returns the sum of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Sum(amounts ...*money.Money) *money.Money {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		a = Add(a, b)
	}

	return a
}
