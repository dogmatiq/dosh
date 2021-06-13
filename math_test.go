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
			"it returns an amount with the absolute magnitude",
			func(dec string, expect string) {
				a := FromString("XYZ", dec).Abs()
				x := FromString("XYZ", expect)
				Expect(a.EqualTo(x)).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.23", "1.23"),
			Entry("negative", "-1.23", "1.23"),
		)
	})

	Describe("func Neg()", func() {
		DescribeTable(
			"it returns an amount with the inverse magnitude",
			func(dec string, expect string) {
				a := FromString("XYZ", dec).Neg()
				x := FromString("XYZ", expect)
				Expect(a.EqualTo(x)).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.23", "-1.23"),
			Entry("negative", "-1.23", "1.23"),
		)
	})

	Describe("func Add()", func() {
		It("returns a + b", func() {
			a := FromString("XYZ", "1.23")
			b := FromString("XYZ", "3.45")
			x := FromString("XYZ", "4.68")
			Expect(a.Add(b).EqualTo(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := FromString("XYZ", "1")
				b := FromString("ABC", "1")
				a.Add(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func Sub()", func() {
		It("returns a - b", func() {
			a := FromString("XYZ", "1.23")
			b := FromString("XYZ", "3.45")
			x := FromString("XYZ", "-2.22")
			Expect(a.Sub(b).EqualTo(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := FromString("XYZ", "1")
				b := FromString("ABC", "1")
				a.Sub(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func MulScalar()", func() {
		It("returns a * b", func() {
			a := FromString("XYZ", "1.23")
			b := decimal.RequireFromString("3.45")
			x := FromString("XYZ", "4.2435")
			Expect(a.MulScalar(b).EqualTo(x)).To(BeTrue())
		})
	})

	Describe("func Div()", func() {
		It("returns a / b", func() {
			a := FromString("XYZ", "1.23")
			b := FromString("XYZ", "0.5")
			x := decimal.RequireFromString("2.46")
			Expect(a.Div(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := FromString("XYZ", "1")
				b := FromString("ABC", "1")
				a.Div(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})

		It("panics when dividing by zero", func() {
			Expect(func() {
				Amount{}.Div(Amount{})
			}).To(PanicWith("decimal division by 0"))
		})
	})

	Describe("func DivScalar()", func() {
		It("returns a / b", func() {
			a := FromString("XYZ", "1.23")
			b := decimal.RequireFromString("0.5")
			x := FromString("XYZ", "2.46")
			Expect(a.DivScalar(b).EqualTo(x)).To(BeTrue())
		})

		It("panics when dividing by zero", func() {
			Expect(func() {
				Amount{}.DivScalar(decimal.Decimal{})
			}).To(PanicWith("decimal division by 0"))
		})
	})

	Describe("func Mod()", func() {
		It("returns a % b", func() {
			a := FromString("XYZ", "1.23")
			b := FromString("XYZ", "0.5")
			x := decimal.RequireFromString("0.23")
			Expect(a.Mod(b).Equal(x)).To(BeTrue())
		})

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := FromString("XYZ", "1")
				b := FromString("ABC", "1")
				a.Mod(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})

		It("panics when dividing by zero", func() {
			Expect(func() {
				Amount{}.Mod(Amount{})
			}).To(PanicWith("decimal division by 0"))
		})
	})

	Describe("func ModScalar()", func() {
		It("returns a % b", func() {
			a := FromString("XYZ", "1.23")
			b := decimal.RequireFromString("0.5")
			x := FromString("XYZ", "0.23")
			Expect(a.ModScalar(b).EqualTo(x)).To(BeTrue())
		})

		It("panics when dividing by zero", func() {
			Expect(func() {
				Amount{}.ModScalar(decimal.Decimal{})
			}).To(PanicWith("decimal division by 0"))
		})
	})
})

var _ = Describe("func Sum()", func() {
	It("returns the sum of all amounts", func() {
		Expect(
			Sum(
				FromString("XYZ", "1"),
				FromString("XYZ", "2"),
				FromString("XYZ", "3"),
			).EqualTo(
				FromString("XYZ", "6"),
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
				FromString("XYZ", "1"),
				FromString("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})

var _ = Describe("func Avg()", func() {
	It("returns the average (mean) of all values", func() {
		Expect(
			Avg(
				FromString("XYZ", "1"),
				FromString("XYZ", "2"),
				FromString("XYZ", "3"),
			).EqualTo(
				FromString("XYZ", "2"),
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
				FromString("XYZ", "1"),
				FromString("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})
