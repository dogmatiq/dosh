package protomoney_test

import (
	. "github.com/dogmatiq/dosh/protomoney"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/type/money"
)

var _ = Describe("func Abs()", func() {
	DescribeTable(
		"it returns an amount with the absolute magnitude",
		func(m, expect *money.Money) {
			Expect(Abs(m)).To(Equal(expect))
		},
		Entry("zero", &money.Money{}, &money.Money{}),
		Entry("positive", &money.Money{Units: 1, Nanos: 230000000}, &money.Money{Units: 1, Nanos: 230000000}),
		Entry("negative", &money.Money{Units: -1, Nanos: -230000000}, &money.Money{Units: 1, Nanos: 230000000}),
	)
})

var _ = Describe("func Neg()", func() {
	DescribeTable(
		"it returns an amount with the inverse magnitude",
		func(m, expect *money.Money) {
			Expect(Neg(m)).To(Equal(expect))
		},
		Entry("zero", &money.Money{}, &money.Money{}),
		Entry("positive", &money.Money{Units: 1, Nanos: 230000000}, &money.Money{Units: -1, Nanos: -230000000}),
		Entry("negative", &money.Money{Units: -1, Nanos: -230000000}, &money.Money{Units: 1, Nanos: 230000000}),
	)
})

var _ = Describe("func Add()", func() {
	It("returns a + b", func() {
		a := &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 230000000}
		b := &money.Money{CurrencyCode: "XYZ", Units: 3, Nanos: 450000000}
		x := &money.Money{CurrencyCode: "XYZ", Units: 4, Nanos: 680000000}
		Expect(Add(a, b)).To(Equal(x))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			a := &money.Money{CurrencyCode: "XYZ"}
			b := &money.Money{CurrencyCode: "ABC"}
			Add(a, b)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})

var _ = Describe("func Sub()", func() {
	It("returns a - b", func() {
		a := &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 230000000}
		b := &money.Money{CurrencyCode: "XYZ", Units: 3, Nanos: 450000000}
		x := &money.Money{CurrencyCode: "XYZ", Units: -2, Nanos: -220000000}
		Expect(Sub(a, b)).To(Equal(x))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			a := &money.Money{CurrencyCode: "XYZ"}
			b := &money.Money{CurrencyCode: "ABC"}
			Sub(a, b)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})

var _ = Describe("func Sum()", func() {
	It("returns the sum of all amounts", func() {
		Expect(
			Sum(
				&money.Money{CurrencyCode: "XYZ", Units: 10, Nanos: 1},
				&money.Money{CurrencyCode: "XYZ", Units: 20, Nanos: 2},
				&money.Money{CurrencyCode: "XYZ", Units: 30, Nanos: 999999999},
			),
		).To(Equal(
			&money.Money{CurrencyCode: "XYZ", Units: 61, Nanos: 2},
		))
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Sum()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Sum(
				&money.Money{CurrencyCode: "XYZ"},
				&money.Money{CurrencyCode: "ABC"},
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})
