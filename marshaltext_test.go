package dosh_test

import (
	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount (text marshaling)", func() {
	Describe("func MarshalText()", func() {
		It("returns a textual representation of the amount", func() {
			a := MustParse("xyz", "10.123")

			data, err := a.MarshalText()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).To(Equal([]byte("XYZ 10.123"))) // note: uppercase
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals an amount from its textual representation", func() {
			var a Amount

			data := []byte("xyz 10.123")
			err := a.UnmarshalText(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

			m := decimal.RequireFromString("10.123")
			Expect(a.Magnitude().Equal(m))
		})

		DescribeTable(
			"it returns an error if the data is invalid",
			func(data string, expect string) {
				var a Amount
				err := a.UnmarshalText([]byte(data))
				Expect(err).To(MatchError(expect))
			},
			Entry("empty", "", "cannot unmarshal amount from text representation: data must have current and magnitude components"),
			Entry("empty currency", " 1.23", "cannot unmarshal amount from text representation: data must have current and magnitude components"),
			Entry("empty magnitude", "XYZ ", "cannot unmarshal amount from text representation: data must have current and magnitude components"),
			Entry("invalid magnitude", "XYZ <invalid>", "cannot unmarshal amount from text representation: can't convert <invalid> to decimal"),
		)
	})
})
