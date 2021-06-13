package dosh

import (
	"fmt"

	"github.com/dogmatiq/dosh/internal/currency"
	"github.com/shopspring/decimal"
)

var (
	// zero is a decimal with a value of 0 (zero).
	zero = decimal.Decimal{}

	// unit is a decimal with a value of 1 (one).
	unit = decimal.NewFromInt(1)
)

// Amount represents an immutable amount of money in a specific currency.
//
// The zero-value represents zero US dollars ($0 USD).
//
// An Amount consists of a currency code and a magnitude.
//
// The currency code is a short string that identifies which currency the amount
// represents. The magnitude is an arbitrary-precision decimal number describing
// the number of units of that currency.
//
// Currency codes must consist of 3 or more uppercase ASCII letters. Any
// operation that accepts a currency code will panic if provided with a code
// that does not meet this criteria. Where possible, currency codes should be an
// ISO-4217 3-letter code. Non-standard currency codes should begin with an "X".
type Amount struct {
	_ [0]func() // prevent comparison with ==

	// cur is the currency code that idenfifies what currency the magnitude is
	// expressed in.
	//
	// An empty string is equivalent to "USD".
	cur string

	// mag is the monetary amount, expressed in whatever currency is specified
	// by the currency field.
	mag decimal.Decimal
}

// Zero returns an Amount with a magnitude of 0 (zero).
//
// c is the currency code that identifies the currency.
func Zero(c string) Amount {
	return FromDecimal(c, zero)
}

// Unit returns an Amount with a magnitude of 1 (one).
//
// c is the currency code that identifies the currency.
func Unit(c string) Amount {
	return FromDecimal(c, unit)
}

// FromDecimal returns an Amount with a decimal magnitude
//
// c is the currency code that identifies the currency.
//
// m is the magnitude of the amount, expressed in the currency specified by c.
func FromDecimal(c string, m decimal.Decimal) Amount {
	if err := currency.ValidateCode(c); err != nil {
		panic(err)
	}

	return Amount{
		cur: c,
		mag: m,
	}
}

// FromInt returns an Amount with an integer magnitude.
//
// c is the currency code that identifies the currency.
//
// m is the magnitude of the amount, expressed in the currency specified by c.
func FromInt(c string, m int) Amount {
	return FromDecimal(c, decimal.NewFromInt(int64(m)))
}

// FromString returns an Amount with a magnitude parsed from a numeric string.
//
// c is the currency code that identifies the currency.
//
// m is the string representation if the magnitude, expressed in the currency
// specified by c. It must use integer, decimal or scientific notation,
// otherwise a panic occurs.
func FromString(c, m string) Amount {
	return FromDecimal(
		c,
		decimal.RequireFromString(m),
	)
}

// TryFromString returns an Amount with a magnitude parsed from a numeric
// string.
//
// c is the currency code that identifies the currency.
//
// m is the string representation if the magnitude, expressed in the currency
// specified by c. It must use integer, decimal or scientific notation,
// otherwise ok is false, and the returned amount is undefined.
func TryFromString(c, m string) (_ Amount, ok bool) {
	d, err := decimal.NewFromString(m)
	return FromDecimal(c, d), err == nil
}

// CurrencyCode returns the currency code for the currency in which the amount
// is specified.
func (a Amount) CurrencyCode() string {
	if len(a.cur) == 0 {
		return "USD"
	}

	return a.cur
}

// Magnitude returns the decimal value of the amount without currency
// information.
func (a Amount) Magnitude() decimal.Decimal {
	return a.mag
}

// String returns a human-readable representation of the amount, including the
// currency code.
func (a Amount) String() string {
	return a.CurrencyCode() + " " + a.mag.String()
}

// GoString returns a string representation of the amount in Go syntax.
func (a Amount) GoString() string {
	if a.mag.IsZero() {
		return fmt.Sprintf(
			"money.Zero(%#v)",
			a.CurrencyCode(),
		)
	}

	if a.mag.Equal(unit) {
		return fmt.Sprintf(
			"money.Unit(%#v)",
			a.CurrencyCode(),
		)
	}

	return fmt.Sprintf(
		"money.FromString(%#v, %#v)",
		a.CurrencyCode(),
		a.mag.String(),
	)
}

// Format implements fmt.Formatter, allowing Amount to be used with fmt.Printf()
// and its variants.
func (a Amount) Format(f fmt.State, verb rune) {
	switch verb {
	case 'f', 'F', // decimal notation
		'e', 'E', // scientific notation
		'g', 'G', // decimal notation, or scientific for large exponents
		'v': // "default" format
		// This is a subset of the verbs are supported by big.Float.Format().
	default:
		fmt.Fprintf(f, "%%!%c(money.Amount=%s)", verb, a.String())
		return
	}

	fmt.Fprintf(f, "%s ", a.CurrencyCode())
	a.mag.BigFloat().Format(f, verb)
}

// assertSameCurrency panics if a and b do not have the same currency.
func assertSameCurrency(a, b Amount) {
	if a.CurrencyCode() != b.CurrencyCode() {
		panic(fmt.Sprintf(
			"can not operate on amounts in differing currencies (%s vs %s)",
			a.CurrencyCode(),
			b.CurrencyCode(),
		))
	}
}
