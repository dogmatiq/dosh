package dosh_test

import (
	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount", func() {
	Describe("func Zero()", func() {
		It("panics if the currency code is empty", func() {
			Expect(func() {
				Zero("")
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func New()", func() {
		It("panics if the currency code is empty", func() {
			Expect(func() {
				New("", decimal.Decimal{})
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func Parse()", func() {
		It("returns an amount parsed from a decimal string", func() {
			a, err := Parse("XYZ", "1.23")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.RequireFromString("1.23")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("returns an error if the input is invalid", func() {
			_, err := Parse("XYZ", "<invalid>")
			Expect(err).Should(HaveOccurred())
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				Parse("", "1.23")
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func MustParse()", func() {
		It("panics if the input is invalid", func() {
			Expect(func() {
				MustParse("XYZ", "<invalid>")
			}).To(Panic())
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				MustParse("", "1.23")
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func CurrencyCode()", func() {
		It("returns USD by default", func() {
			var a Amount
			Expect(a.CurrencyCode()).To(Equal("USD"))
		})

		It("returns the provided code in upper case", func() {
			a := Zero("xyz")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))
		})
	})

	Describe("func Magnitude()", func() {
		It("returns zero by default", func() {
			var a Amount
			Expect(a.Magnitude().Equal(decimal.Decimal{})).To(BeTrue())
		})

		It("returns the amount's magnitude", func() {
			m := decimal.RequireFromString("1.23")
			a := New("XYZ", m)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})
	})

	Describe("func String()", func() {
		It("returns a string representation of the amount", func() {
			a := MustParse("XYZ", "10.123")
			Expect(a.String()).To(Equal("XYZ 10.123"))
		})
	})
})
