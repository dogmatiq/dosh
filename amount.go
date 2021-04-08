package dosh

import (
	"strings"

	"github.com/shopspring/decimal"
)

// Amount represents an amount of money in a specific currency.
//
// The zero-value represents zero US dollars (0 USD).
type Amount struct {
	currency  string
	magnitude decimal.Decimal
}

// New returns an Amount with a specific currency and magnitude.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217, however it MAY be any
// string when used to identify a non-standard currency.
func New(c string, magnitude decimal.Decimal) Amount {
	if c == "" {
		panic("currency code must not be empty")
	}

	return Amount{
		strings.ToUpper(c),
		magnitude,
	}
}

// Zero returns an Amount with a magnitude of zero in a specific currency.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217, however it MAY be any
// string when used to identify a non-standard currency.
func Zero(c string) Amount {
	return New(c, decimal.Decimal{})
}

// Parse returns a new Amount with a magnitude parsed from a decimal string or
// panics if unable to do so.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217, however it MAY be any
// string when used to identify a non-standard currency.
func Parse(c string, magnitude string) (Amount, error) {
	m, err := decimal.NewFromString(magnitude)
	if err != nil {
		return Amount{}, err
	}

	return New(c, m), nil
}

// MustParse returns a new AMount with a magnitude parsed from a decimal string
// or panics if unable to do so.
//
// c is the currency code that identifies the currency. It SHOULD be an ISO-4217
// 3-letter code if the currency is defined by ISO-4217, however it MAY be any
// string when used to identify a non-standard currency.
func MustParse(c string, magnitude string) Amount {
	return New(
		c,
		decimal.RequireFromString(magnitude),
	)
}

// CurrencyCode returns the currency code for the currency in which the amount
// is specified.
//
// The code SHOULD be an ISO-4217 3-letter code if the currency is defined by
// ISO-4217, however the code MAY identify a non-standard currency.
//
// The returned currency code is always uppercase.
func (a Amount) CurrencyCode() string {
	if a.currency == "" {
		return "USD"
	}

	return a.currency
}

// Magnitude returns the decimal value of the amount.
func (a Amount) Magnitude() decimal.Decimal {
	return a.magnitude
}

// func (a Amount) Format(f fmt.State, verb rune) {
// 	panic("not implemented")
// }

// String returns a human-readable representation of the amount, including the
// currency code.
func (a Amount) String() string {
	return a.CurrencyCode() + " " + a.magnitude.String()
}
