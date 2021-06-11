package dosh

import (
	"fmt"
	"strings"

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
// An Amount consists of a currency, specified by a currency code, and a
// "magnitude", which is an arbitrary-precision decimal number describing the
// number of units of that currency.
//
// The zero-value represents zero US dollars (0 USD).
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

// New returns an Amount with a specific currency and magnitude.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
//
// m is the magnitude of the amount, expressed in the currency specified by c.
func New(c string, m decimal.Decimal) Amount {
	if c == "" {
		panic("currency code must not be empty")
	}

	return Amount{
		cur: strings.ToUpper(c),
		mag: m,
	}
}

// Zero returns an Amount with a magnitude of 0 (zero) in a specific currency.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
func Zero(c string) Amount {
	return New(c, zero)
}

// Unit returns an Amount with a magnitude of 1 (one) in a specific currency.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
func Unit(c string) Amount {
	return New(c, unit)
}

// Int returns an Amount with an integer magnitude in a specific currency.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
//
// m is the magnitude of the amount, expressed in the currency specified by c.
func Int(c string, m int) Amount {
	return New(c, decimal.NewFromInt(int64(m)))
}

// Parse returns a new Amount with a magnitude parsed from a decimal string.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
//
// m is the string representation of an integer or decimal number, expressed in
// the currency specified by c.
func Parse(c, m string) (Amount, error) {
	d, err := decimal.NewFromString(m)
	if err != nil {
		return Amount{}, err
	}

	return New(c, d), nil
}

// MustParse returns a new Amount with a magnitude parsed from a decimal string
// or panics if unable to do so.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217. Non-standard currency
// codes SHOULD begin with an "X". All values are converted to uppercase.
//
// m is the string representation of an integer or decimal number, expressed in
// the currency specified by c.
func MustParse(c, m string) Amount {
	return New(
		c,
		decimal.RequireFromString(m),
	)
}

// CurrencyCode returns the currency code for the currency in which the amount
// is specified.
//
// The returned currency code is always uppercase.
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
		"money.MustParse(%#v, %#v)",
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
