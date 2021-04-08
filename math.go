package dosh

import "github.com/shopspring/decimal"

// Abs returns the absolute value of this amount.
//
// That is, if a < 0 it returns -a. Otherwise it returns a unchanged.
func (a Amount) Abs() Amount {
	a.mag = a.mag.Abs()
	return a
}

// Neg returns -a.
func (a Amount) Neg() Amount {
	a.mag = a.mag.Neg()
	return a
}

// Add returns a + b.
//
// It panics if a and b do not use the same currency.
func (a Amount) Add(b Amount) Amount {
	assertSameCurrency(a, b)
	a.mag = a.mag.Add(b.mag)
	return a
}

// Sub returns a - b.
//
// It panics if a and b do not use the same currency.
func (a Amount) Sub(b Amount) Amount {
	assertSameCurrency(a, b)
	a.mag = a.mag.Sub(b.mag)
	return a
}

// Mul returns a * b, where b is a scalar decimal value.
func (a Amount) MulScalar(b decimal.Decimal) Amount {
	a.mag = a.mag.Mul(b)
	return a
}

// Div returns a / b.
//
// It panics if a and b do not use the same currency.
//
// To divide by a scalar decimal value, use DivScalar() instead.
func (a Amount) Div(b Amount) decimal.Decimal {
	assertSameCurrency(a, b)
	return a.mag.Div(b.mag)
}

// DivScalar returns a / b, where b is a scalar value.
//
// To divide a by another Amount, use Div() instead.
func (a Amount) DivScalar(b decimal.Decimal) Amount {
	a.mag = a.mag.Div(b)
	return a
}

// Mod returns a % b.
//
// It panics if a and b do not use the same currency.
//
// To find the remainder of dividing by a scalar decimal value, use ModScalar()
// instead.
func (a Amount) Mod(b Amount) decimal.Decimal {
	assertSameCurrency(a, b)
	return a.mag.Mod(b.mag)
}

// ModScalar returns a % b, where b is a scalar decimal value.
//
// To find the remainder of dividing by another Amount, use Mod() instead.
func (a Amount) ModScalar(b decimal.Decimal) Amount {
	a.mag = a.mag.Mod(b)
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
	return sum.DivScalar(n)
}
