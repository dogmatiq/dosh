package dosh

import (
	"fmt"
)

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
