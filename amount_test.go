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
	Describe("func New()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			m := decimal.NewFromInt(123)
			a := New("xyz", m)
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase
			Expect(a.Magnitude().Equal(m))
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				New("", decimal.Decimal{})
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func Zero()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := Zero("xyz")
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase
			Expect(a.Magnitude().IsZero()).To(BeTrue())
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				Zero("")
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func Unit()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := Unit("xyz")
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

			m := decimal.NewFromInt(1)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				Zero("")
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func Int()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a := Int("xyz", 123)
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

			m := decimal.NewFromInt(123)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

		It("panics if the currency code is empty", func() {
			Expect(func() {
				Int("", 0)
			}).To(PanicWith("currency code must not be empty"))
		})
	})

	Describe("func Parse()", func() {
		It("returns an amount with the correct currency code and magnitude", func() {
			a, err := Parse("xyz", "1.23")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

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
		It("returns an amount with the correct currency code and magnitude", func() {
			a := MustParse("xyz", "1.23")
			Expect(a.CurrencyCode()).To(Equal("XYZ")) // note: uppercase

			m := decimal.RequireFromString("1.23")
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})

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
		It("returns USD when called on a zero-value amount", func() {
			var a Amount
			Expect(a.CurrencyCode()).To(Equal("USD"))
		})

		It("returns the provided code in upper case", func() {
			a := Zero("xyz")
			Expect(a.CurrencyCode()).To(Equal("XYZ"))
		})
	})

	Describe("func Magnitude()", func() {
		It("returns zero when called on a zero-value amount", func() {
			var a Amount
			Expect(a.Magnitude().IsZero()).To(BeTrue())
		})

		It("returns the amount's magnitude", func() {
			m := decimal.RequireFromString("1.23")
			a := New("XYZ", m)
			Expect(a.Magnitude().Equal(m)).To(BeTrue())
		})
	})

	Describe("func String()", func() {
		It("returns a string representation of the amount", func() {
			a := MustParse("xyz", "10.123")
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
			Entry("zero-magnitude", Zero("xyz"), `money.Zero("XYZ")`),
			Entry("unit-magnitude", Unit("xyz"), `money.Unit("XYZ")`),
			Entry("other", MustParse("xyz", "1.23"), `money.MustParse("XYZ", "1.23")`),
		)
	})

	Describe("func Format()", func() {
		It("returns a formatted representation of the amount", func() {
			a := MustParse("xyz", "10.129")
			s := fmt.Sprintf("%0.2f", a)
			Expect(s).To(Equal("XYZ 10.13"))
		})

		It("returns a descriptive string if used with an unsupported verb", func() {
			a := MustParse("xyz", "10.129")
			s := fmt.Sprintf("%d", a)
			Expect(s).To(Equal("%!d(money.Amount=XYZ 10.129)"))
		})
	})
})
