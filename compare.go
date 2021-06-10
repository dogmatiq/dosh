package dosh

// IsZero returns true if the amount has a magnitude of zero.
func (a Amount) IsZero() bool {
	return a.mag.IsZero()
}

// IsNegative returns true if the amount has a negative magnitude.
func (a Amount) IsNegative() bool {
	return a.mag.IsNegative()
}

// IsPositive returns true if the amount has a positive magnitude.
func (a Amount) IsPositive() bool {
	return a.mag.IsPositive()
}

// Cmp compares a to b and returns a C-style comparison result.
//
// It panics if a and b do not use the same currency.
//
// If a < b then c is negative.
// If a > b then c is positive.
// Otherwise; a == b and c is zero.
func (a Amount) Cmp(b Amount) (c int) {
	assertSameCurrency(a, b)
	return a.mag.Cmp(b.mag)
}

// EqualTo returns true if a and b have the same magnitude.
//
// It panics if a and b do not use the same currency.
//
// To check equality between two amounts that may have differing currencies, use
// IdenticalTo() instead.
func (a Amount) EqualTo(b Amount) bool {
	assertSameCurrency(a, b)
	return a.mag.Equal(b.mag)
}

// Identical returns true if a and b use the same currency and have the same
// magnitude.
//
// For general comparisons that are expected to be in the same currency, use
// EqualTo() instead.
func (a Amount) IdenticalTo(b Amount) bool {
	if a.CurrencyCode() != b.CurrencyCode() {
		return false
	}

	return a.mag.Equal(b.mag)
}

// LessThan returns true if a < b.
//
// It panics if a and b do not use the same currency.
func (a Amount) LessThan(b Amount) bool {
	assertSameCurrency(a, b)
	return a.mag.LessThan(b.mag)
}

// LessThanOrEqual returns true if a <= b.
//
// It panics if a and b do not use the same currency.
func (a Amount) LessThanOrEqualTo(b Amount) bool {
	assertSameCurrency(a, b)
	return a.mag.LessThanOrEqual(b.mag)
}

// GreaterThan returns true if a > b.
//
// It panics if a and b do not use the same currency.
func (a Amount) GreaterThan(b Amount) bool {
	assertSameCurrency(a, b)
	return a.mag.GreaterThan(b.mag)
}

// GreaterThanOrEqual returns true if a >= b.
//
// It panics if a and b do not use the same currency.
func (a Amount) GreaterThanOrEqualTo(b Amount) bool {
	assertSameCurrency(a, b)
	return a.mag.GreaterThanOrEqual(b.mag)
}

// LexicallyLessThan returns true if a should appear before b in a sorted list.
//
// There is no requirement that a and b use the same currency.
func (a Amount) LexicallyLessThan(b Amount) bool {
	if a.CurrencyCode() == b.CurrencyCode() {
		return a.mag.LessThan(b.mag)
	}

	return a.CurrencyCode() < b.CurrencyCode()
}

// Min returns the smallest of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Min(amounts ...Amount) Amount {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		assertSameCurrency(a, b)
		if b.LessThan(a) {
			a = b
		}
	}

	return a
}

// Max returns the largest of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Max(amounts ...Amount) Amount {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		assertSameCurrency(a, b)
		if b.GreaterThan(a) {
			a = b
		}
	}

	return a
}
