package dosh

import "github.com/shopspring/decimal"

// Abs returns the absolute value of this amount.
//
// That is, if a < 0 it returns -a. Otherwise it returns a unchanged.
func (a Amount) Abs() Amount {
	a.magnitude = a.magnitude.Abs()
	return a
}

// Neg returns -a.
func (a Amount) Neg() Amount {
	a.magnitude = a.magnitude.Neg()
	return a
}

// Add returns a + b.
//
// It panics if a and b do not use the same currency.
func (a Amount) Add(b Amount) Amount {
	assertSameCurrency(a, b)
	a.magnitude = a.magnitude.Add(b.magnitude)
	return a
}

// Sub returns a - b.
//
// It panics if a and b do not use the same currency.
func (a Amount) Sub(b Amount) Amount {
	assertSameCurrency(a, b)
	a.magnitude = a.magnitude.Sub(b.magnitude)
	return a
}

// Mul returns a * b, where b is a magnitude without any specific currency.
func (a Amount) MulM(b decimal.Decimal) Amount {
	a.magnitude = a.magnitude.Mul(b)
	return a
}

// Div returns a / b.
//
// It panics if a and b do not use the same currency.
//
// To divide by a magnitude without any specific currency, use DivM() instead.
func (a Amount) Div(b Amount) decimal.Decimal {
	assertSameCurrency(a, b)
	return a.magnitude.Div(b.magnitude)
}

// DivM returns a / b, where b is a magnitude without any specific currency.
//
// To divide a by another Amount, use Div() instead.
func (a Amount) DivM(b decimal.Decimal) Amount {
	a.magnitude = a.magnitude.Div(b)
	return a
}

// Mod returns a % b.
//
// It panics if a and b do not use the same currency.
//
// To find the remainder of dividing by a magnitude without any specific
// currency, use ModM() instead.
func (a Amount) Mod(b Amount) decimal.Decimal {
	assertSameCurrency(a, b)
	return a.magnitude.Mod(b.magnitude)
}

// ModM returns a % b, where b is a magnitude without any specific currency.
//
// To find the remainder of dividing by another Amount, use Mod() instead.
func (a Amount) ModM(b decimal.Decimal) Amount {
	a.magnitude = a.magnitude.Mod(b)
	return a
}

// Sum returns the sum of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Sum(amounts ...Amount) Amount {
	if len(amounts) == 0 {
		panic("at least one amount must be provided")
	}

	a := amounts[0]
	for _, b := range amounts[1:] {
		a = a.Add(b)
	}

	return a
}

// Avg returns the mean of the given amounts.
//
// It panics if amounts is empty, or if the amounts do not use the same
// currency.
func Avg(amounts ...Amount) Amount {
	sum := Sum(amounts...)
	n := decimal.NewFromInt(int64(len(amounts)))
	return sum.DivM(n)
}
