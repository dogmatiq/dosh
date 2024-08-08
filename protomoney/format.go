package protomoney

import (
	"fmt"
	"math/big"

	"google.golang.org/genproto/googleapis/type/money"
)

// Fmt wraps a money value in a formatter that allows it to be formatted using
// standard fmt.Printf() verbs.
func Fmt(m *money.Money) fmt.Formatter {
	return formatter{m}
}

type formatter struct {
	value *money.Money
}

func (f formatter) Format(st fmt.State, verb rune) {
	units, nanos := normalizeComponents(f.value)

	switch verb {
	case 'f', 'F', // decimal notation
		'e', 'E', // scientific notation
		'g', 'G', // decimal notation, or scientific for large exponents
		'v': // "default" format
		// This is a subset of the verbs are supported by big.Float.Format().

		s := fmt.Sprintf("%d.%09d", units, nanos)
		v := &big.Float{}
		v.SetString(s)

		fmt.Fprintf(st, "%s ", f.value.GetCurrencyCode())
		v.Format(st, verb)

	default:
		fmt.Fprintf(st, "%%!%c(*money.Money=%s)", verb, f.value.String())
	}
}
