package dosh_test

import (
	"math"

	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount (JSON marshaling)", func() {
	Describe("func MarshalJSON()", func() {
		It("returns the JSON representation of the amount", func() {
			a := MustParse("xyz", "10.123")

			data, err := a.MarshalJSON()
			Expect(err).ShouldNot(HaveOccurred())

			// note: uppercase currency, and units represented as a string
			Expect(data).To(MatchJSON([]byte(`{"currency_code":"XYZ","units":"10","nanos":123000000}`)))
		})

		DescribeTable(
			"it returns an error if the amount can not be marshaled",
			func(a Amount, expect string) {
				_, err := a.MarshalJSON()
				Expect(err).To(MatchError(expect))
			},
			Entry(
				"integer component of the magnitude overflows an int64",
				Int("xyz", math.MaxInt64).Add(Unit("xyz")),
				"cannot marshal amount to JSON representation: magnitude's integer component overflows int64",
			),
			Entry(
				"fractional component of the magnitude requires more precision that available",
				MustParse("xyz", "0.0123456789"),
				"cannot marshal amount to JSON representation: magnitude's fractional component has too many decimal places",
			),
		)
	})

	Describe("func UnmarshalJSON()", func() {
		It("unmarshals an amount from its JSON representation", func() {
			var a Amount

			err := a.UnmarshalJSON([]byte(`{"currency_code":"xyz","units":"10","nanos":123000000}`))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

			m := decimal.RequireFromString("10.123")
			Expect(a.Magnitude().Equal(m))
		})

		DescribeTable(
			"it returns an error if the JSON message is invalid",
			func(data string, expect string) {
				var a Amount
				err := a.UnmarshalJSON([]byte(data))
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp(expect))
			},
			Entry(
				"malformed JSON",
				`<invalid>`,
				// Note, protocol buffers package randomly emits different
				// output to avoid code that checks errors by string value,
				// which is fair enough but makes simple tests like this
				// incredibly frustrating. This is the reason for the use of
				// regular expressions.
				"cannot unmarshal amount from JSON representation: proto:.+syntax error",
			),
			Entry(
				"empty currency",
				`{"units":"0","nanos":0}`,
				"cannot unmarshal amount from JSON representation: currency code must not be empty",
			),
			Entry(
				"units positive, nanos negative",
				`{"currency_code": "XYZ", "units": "1", "nanos": -1}`,
				"cannot unmarshal amount from JSON representation: units and nanos components must have the same sign",
			),
			Entry(
				"units negative, nanos positive",
				`{"currency_code": "XYZ", "units": "-1", "nanos": 1}`,
				"cannot unmarshal amount from JSON representation: units and nanos components must have the same sign",
			),
		)
	})
})
