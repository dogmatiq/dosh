package dosh_test

import (
	"fmt"

	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount", func() {
	Describe("func Zero()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := Zero("XYZ")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))
			Expect(a.Magnitude().IsZero()).To(BeTrue())
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				Zero("X")
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func Unit()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := Unit("XYZ")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.NewFromInt(1)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				Unit("X")
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func FromDecimal()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			m := decimal.NewFromInt(123)
			a := FromDecimal("XYZ", m)
			Expect(a.CurrencyCode()).To(Equal("XYZ"))
			Expect(a.Magnitude().Equal(m))
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				FromDecimal("X", decimal.Decimal{})
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func FromInt()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := FromInt("XYZ", 123)
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.NewFromInt(123)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				FromInt("X", 0)
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func FromString()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			By("parsing an integer")

			a := FromString("XYZ", "1")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.RequireFromString("1")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())

			By("parsing a decimal")

			a = FromString("XYZ", "1.23")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m = decimal.RequireFromString("1.23")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())

			By("parsing a decimal with trailing zeroes")

			a = FromString("XYZ", "1000.00")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m = decimal.RequireFromString("1000.00")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())

			By("parsing scientific notation")

			a = FromString("XYZ", "123e-2")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m = decimal.RequireFromString("1.23")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("panics if the input is invalid", func() {
			Expect(func() {
				FromString("XYZ", "<invalid>")
			}).To(Panic())
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				FromString("X", "1.23")
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func TryFromString()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a, ok := TryFromString("XYZ", "1.23")
			Expect(ok).To(BeTrue())
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.RequireFromString("1.23")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("returns false if the input is invalid", func() {
			_, ok := TryFromString("XYZ", "<invalid>")
			Expect(ok).To(BeFalse())
		})

		It("panics if the currency code is invalid", func() {
			Expect(func() {
				TryFromString("X", "1.23")
			}).To(PanicWith(MatchError("currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters")))
		})
	})

	Describe("func CurrencyCode()", func() {
		It("returns the currency code", func() {
			a := Zero("XYZ")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))
		})

		It("returns USD when called on a zero-value amount", func() {
			var a Amount
			Expect(a.CurrencyCode()).To(Equal("USD"))
		})
	})

	Describe("func Magnitude()", func() {
		It("returns zero when called on a zero-value amount", func() {
			var a Amount
			Expect(a.Magnitude().IsZero()).To(BeTrue())
		})

		It("returns the amount's magnitude", func() {
			m := decimal.RequireFromString("1.23")
			a := FromDecimal("XYZ", m)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})
	})

	Describe("func String()", func() {
		It("returns a string representation of the amount", func() {
			a := FromString("XYZ", "10.123")
			Expect(a.String()).To(Equal("XYZ 10.123"))
		})
	})

	Describe("func GoString()", func() {
		DescribeTable(
			"it returns a string representation of the amount in Go syntax",
			func(a Amount, expect string) {
				Expect(a.GoString()).To(Equal(expect))
			},
			Entry("zero-value", Amount{}, `money.Zero("USD")`),
			Entry("zero-magnitude", Zero("XYZ"), `money.Zero("XYZ")`),
			Entry("unit-magnitude", Unit("XYZ"), `money.Unit("XYZ")`),
			Entry("other", FromString("XYZ", "1.23"), `money.FromString("XYZ", "1.23")`),
		)
	})

	Describe("func Format()", func() {
		It("returns a formatted representation of the amount", func() {
			a := FromString("XYZ", "10.129")
			s := fmt.Sprintf("%0.2f", a)
			Expect(s).To(Equal("XYZ 10.13"))
		})

		It("returns a descriptive string if used with an unsupported verb", func() {
			a := FromString("XYZ", "10.129")
			s := fmt.Sprintf("%d", a)
			Expect(s).To(Equal("%!d(money.Amount=XYZ 10.129)"))
		})
	})
})
