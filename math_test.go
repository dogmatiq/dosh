package dosh_test

import (
	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount (math methods)", func() {
	Describe("func Abs()", func() {
		DescribeTable(
			"returns an amount with the absolute magnitude",
			func(dec string, expect string) {
				a := MustParse("XYZ", dec).Abs()
				x := MustParse("XYZ", expect)
				Expect(a.Equal(x)).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.23", "1.23"),
			Entry("negative", "-1.23", "1.23"),
		)
	})

	Describe("func Neg()", func() {
		DescribeTable(
			"returns an amount with the inverse magnitude",
			func(dec string, expect string) {
				a := MustParse("XYZ", dec).Neg()
				x := MustParse("XYZ", expect)
				Expect(a.Equal(x)).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.23", "-1.23"),
			Entry("negative", "-1.23", "1.23"),
		)
	})

	Describe("func Add()", func() {
		It("returns a + b", func() {
			a := MustParse("XYZ", "1.23")
			b := MustParse("XYZ", "3.45")
			x := MustParse("XYZ", "4.68")
			Expect(a.Add(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.Add(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func Sub()", func() {
		It("returns a - b", func() {
			a := MustParse("XYZ", "1.23")
			b := MustParse("XYZ", "3.45")
			x := MustParse("XYZ", "-2.22")
			Expect(a.Sub(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.Sub(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func MulM()", func() {
		It("returns a * b", func() {
			a := MustParse("XYZ", "1.23")
			b := decimal.RequireFromString("3.45")
			x := MustParse("XYZ", "4.2435")
			Expect(a.MulM(b).Equal(x)).To(BeTrue())
		})
	})

	Describe("func Div()", func() {
		It("returns a / b", func() {
			a := MustParse("XYZ", "1.23")
			b := MustParse("XYZ", "0.5")
			x := decimal.RequireFromString("2.46")
			Expect(a.Div(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.Div(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func DivM()", func() {
		It("returns a / b", func() {
			a := MustParse("XYZ", "1.23")
			b := decimal.RequireFromString("0.5")
			x := MustParse("XYZ", "2.46")
			Expect(a.DivM(b).Equal(x)).To(BeTrue())
		})
	})

	Describe("func Mod()", func() {
		It("returns a % b", func() {
			a := MustParse("XYZ", "1.23")
			b := MustParse("XYZ", "0.5")
			x := decimal.RequireFromString("0.23")
			Expect(a.Mod(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.Mod(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func ModM()", func() {
		It("returns a % b", func() {
			a := MustParse("XYZ", "1.23")
			b := decimal.RequireFromString("0.5")
			x := MustParse("XYZ", "0.23")
			Expect(a.ModM(b).Equal(x)).To(BeTrue())
		})
	})
})

var _ = Describe("func Sum()", func() {
	It("returns the sum of all amounts", func() {
		Expect(
			Sum(
				MustParse("XYZ", "1"),
				MustParse("XYZ", "2"),
				MustParse("XYZ", "3"),
			).Equal(
				MustParse("XYZ", "6"),
			),
		).To(BeTrue())
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Sum()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Sum(
				MustParse("XYZ", "1"),
				MustParse("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})

var _ = Describe("func Avg()", func() {
	It("returns the average (mean) of all values", func() {
		Expect(
			Avg(
				MustParse("XYZ", "1"),
				MustParse("XYZ", "2"),
				MustParse("XYZ", "3"),
			).Equal(
				MustParse("XYZ", "2"),
			),
		).To(BeTrue())
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Avg()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Avg(
				MustParse("XYZ", "1"),
				MustParse("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})
