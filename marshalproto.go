package dosh

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
)

var (
	unit         = decimal.NewFromInt(1)
	nanosPerUnit = decimal.NewFromInt(1_000_000_000)
)

// MarshalProto mashals an amount to its protocol buffers representation.
func (a Amount) MarshalProto() (*money.Money, error) {
	if !a.magnitude.BigInt().IsInt64() {
		return nil, errors.New("cannot marshal amount to protocol buffers representation: integer component is too large")
	}

	// Isolate the fractional part and multiply it by nanosPerUnit to work out
	// how many nano units we have.
	nanos := a.magnitude.
		Mod(unit).
		Mul(nanosPerUnit)

	// If after doing so nanos still contains a fractional component then
	// the decimal has more decimal places than can be represented using
	// nano units.
	if !nanos.Mod(unit).Equal(decimal.Zero) {
		return nil, errors.New("cannot marshal amount to protocol buffers representation: fractional component has too many decimal places")
	}

	return &money.Money{
		CurrencyCode: a.CurrencyCode(),
		Units:        a.magnitude.IntPart(),
		Nanos:        int32(nanos.IntPart()),
	}, nil
}

// UnmarshalProto unmarshals an amount from its protocol buffers representation.
func (a *Amount) UnmarshalProto(pb *money.Money) error {
	c := strings.TrimSpace(pb.GetCurrencyCode())
	if c == "" {
		return errors.New("cannot unmarshal amount from text representation: currency code must not be empty")
	}

	units := pb.GetUnits()
	nanos := int64(pb.GetNanos())

	if (units > 0 && nanos < 0) || (units < 0 && nanos > 0) {
		return errors.New("cannot unmarshal amount from text representation: units and nanos components must have the same sign")
	}

	// normalize any overflowing nanos component into units.
	units += int64(nanos) / nanosPerUnit.IntPart()
	nanos %= nanosPerUnit.IntPart()

	a.currency = strings.ToUpper(c)
	a.magnitude = decimal.NewFromInt(units).Add(
		decimal.NewFromInt(nanos).Div(nanosPerUnit),
	)

	return nil
}
