package dosh

import (
	"errors"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
)

// nanosPerUnit is the number of "nano units" in each unit.
var nanosPerUnit = decimal.NewFromInt(1_000_000_000)

// MarshalProto mashals an amount to its protocol buffers representation.
func (a Amount) MarshalProto() (*money.Money, error) {
	pb, err := a.marshalProto()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal amount to protocol buffers representation: %w", err)
	}

	return pb, nil
}

// UnmarshalProto unmarshals an amount from its protocol buffers representation.
//
// NOTE: For consistency with other UnmarshalXXX() methods, this method mutates
// the internals of a, violating Amount's immutability guarantee.
func (a *Amount) UnmarshalProto(pb *money.Money) error {
	if err := a.unmarshalProto(pb); err != nil {
		return fmt.Errorf("cannot unmarshal amount from protocol buffers representation: %w", err)
	}

	return nil
}

// marshalProto mashals an amount to its protocol buffers representation,
// without providing any protocol-buffer-specific error information, allowing it
// to be used by both MarshalJSON() and MarshalProto().
func (a Amount) marshalProto() (*money.Money, error) {
	if !a.mag.BigInt().IsInt64() {
		return nil, errors.New("magnitude's integer component overflows int64")
	}

	// Isolate the fractional part and multiply it by nanosPerUnit to work out
	// how many nano units we have.
	nanos := a.mag.
		Mod(unit).
		Mul(nanosPerUnit)

	// If after doing so nanos still contains a fractional component then
	// the decimal has more decimal places than can be represented using
	// nano units.
	if !nanos.Mod(unit).Equal(decimal.Zero) {
		return nil, errors.New("magnitude's fractional component has too many decimal places")
	}

	return &money.Money{
		CurrencyCode: a.CurrencyCode(),
		Units:        a.mag.IntPart(),
		Nanos:        int32(nanos.IntPart()),
	}, nil
}

// unmarshalProto unmashals an amount from its protocol buffers representation,
// without providing any protocol-buffer-specific error information, allowing it
// to be used by both UnmarshalJSON() and UnmarshalProto().
func (a *Amount) unmarshalProto(pb *money.Money) error {
	c := strings.TrimSpace(pb.GetCurrencyCode())
	if c == "" {
		return errors.New("currency code must not be empty")
	}

	units := pb.GetUnits()
	nanos := int64(pb.GetNanos())

	if (units > 0 && nanos < 0) || (units < 0 && nanos > 0) {
		return errors.New("units and nanos components must have the same sign")
	}

	// normalize any overflowing nanos component into units.
	units += int64(nanos) / nanosPerUnit.IntPart()
	nanos %= nanosPerUnit.IntPart()

	a.cur = []byte(strings.ToUpper(c))
	a.mag = decimal.NewFromInt(units).Add(
		decimal.NewFromInt(nanos).Div(nanosPerUnit),
	)

	return nil
}
