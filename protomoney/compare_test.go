package protomoney_test

import (
	"sort"

	. "github.com/dogmatiq/dosh/protomoney"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/type/money"
)

var _ = Describe("func IsZero()", func() {
	DescribeTable(
		"it returns true if the amount has a magnitude of zero",
		func(m *money.Money, expect bool) {
			Expect(IsZero(m)).To(Equal(expect))
		},
		Entry("zero", &money.Money{}, true),
		Entry("positive units", &money.Money{Units: 1}, false),
		Entry("positive nanos", &money.Money{Nanos: 1}, false),
		Entry("positive units & nanos", &money.Money{Units: 1, Nanos: 1}, false),
		Entry("negative units", &money.Money{Units: -1}, false),
		Entry("negative nanos", &money.Money{Nanos: -1}, false),
		Entry("negative units & nanos", &money.Money{Units: -1, Nanos: -1}, false),
	)
})

var _ = Describe("func IsPositive()", func() {
	DescribeTable(
		"it returns true if the amount has a positive magnitude",
		func(m *money.Money, expect bool) {
			Expect(IsPositive(m)).To(Equal(expect))
		},
		Entry("zero", &money.Money{}, false),
		Entry("positive units", &money.Money{Units: 1}, true),
		Entry("positive nanos", &money.Money{Nanos: 1}, true),
		Entry("positive units & nanos", &money.Money{Units: 1, Nanos: 1}, true),
		Entry("negative units", &money.Money{Units: -1}, false),
		Entry("negative nanos", &money.Money{Nanos: -1}, false),
		Entry("negative units & nanos", &money.Money{Units: -1, Nanos: -1}, false),
	)
})

var _ = Describe("func IsNegative()", func() {
	DescribeTable(
		"it returns true if the amount has a negative magnitude",
		func(m *money.Money, expect bool) {
			Expect(IsNegative(m)).To(Equal(expect))
		},
		Entry("zero", &money.Money{}, false),
		Entry("positive units", &money.Money{Units: 1}, false),
		Entry("positive nanos", &money.Money{Nanos: 1}, false),
		Entry("positive units & nanos", &money.Money{Units: 1, Nanos: 1}, false),
		Entry("negative units", &money.Money{Units: -1}, true),
		Entry("negative nanos", &money.Money{Nanos: -1}, true),
		Entry("negative units & nanos", &money.Money{Units: -1, Nanos: -1}, true),
	)
})

var _ = Context("comparison functions", func() {
	vectors := []TableEntry{
		Entry("equal, zero",
			&money.Money{},
			&money.Money{},
			0,
		),
		Entry("equal, positive units",
			&money.Money{Units: 1},
			&money.Money{Units: 1},
			0,
		),
		Entry("equal, positive nanos",
			&money.Money{Nanos: 1},
			&money.Money{Nanos: 1},
			0,
		),
		Entry("equal, positive units & nanos",
			&money.Money{Units: 1, Nanos: 1},
			&money.Money{Units: 1, Nanos: 1},
			0,
		),
		Entry("equal, positive, denormalized",
			&money.Money{Units: 2, Nanos: 500000000},
			&money.Money{Units: 1, Nanos: 1500000000},
			0,
		),
		Entry("equal, negative units",
			&money.Money{Units: -1},
			&money.Money{Units: -1},
			0,
		),
		Entry("equal, negative nanos",
			&money.Money{Nanos: -1},
			&money.Money{Nanos: -1},
			0,
		),
		Entry("equal, negative units & nanos",
			&money.Money{Units: -1, Nanos: -1},
			&money.Money{Units: -1, Nanos: -1},
			0,
		),
		Entry("equal, negative, denormalized",
			&money.Money{Units: -2, Nanos: -500000000},
			&money.Money{Units: -1, Nanos: -1500000000},
			0,
		),

		Entry("less, zero",
			&money.Money{},
			&money.Money{Units: 1}, -1,
		),
		Entry("less, positive units",
			&money.Money{Units: +1},
			&money.Money{Units: +2}, -1,
		),
		Entry("less, positive nanos",
			&money.Money{Nanos: +1},
			&money.Money{Nanos: +2}, -1,
		),
		Entry("less, positive units & nanos",
			&money.Money{Units: +1, Nanos: +1},
			&money.Money{Units: +1, Nanos: +2}, -1,
		),
		Entry("less, negative units",
			&money.Money{Units: -2},
			&money.Money{Units: -1}, -1,
		),
		Entry("less, negative nanos",
			&money.Money{Nanos: -2},
			&money.Money{Nanos: +1}, -1,
		),
		Entry("less, negative units & nanos",
			&money.Money{Units: -1, Nanos: -2},
			&money.Money{Units: -1, Nanos: -1}, -1,
		),
		Entry("less, denormalized",
			&money.Money{Units: 0, Nanos: 1500000000},
			&money.Money{Units: 1, Nanos: 1500000000}, -1,
		),
	}

	Describe("func Cmp()", func() {
		DescribeTable(
			"it returns a C-style comparison result",
			func(a, b *money.Money, expect int) {
				x := Cmp(a, b)
				y := Cmp(b, a)

				if expect == 0 {
					Expect(x).To(BeZero())
					Expect(y).To(BeZero())
				} else {
					Expect(x < 0).To(Equal(expect < 0))
					Expect(y < 0).To(Equal(expect > 0))
				}
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				Cmp(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})

		It("panics if either of the operands units/nanos signs disagree", func() {
			a := &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: -1}
			b := &money.Money{CurrencyCode: "XYZ"}

			Expect(func() {
				Cmp(a, b)
			}).To(PanicWith(MatchError("sign of units component (1) does not agree with sign of nanos component (-1)")))

			Expect(func() {
				Cmp(b, a)
			}).To(PanicWith(MatchError("sign of units component (1) does not agree with sign of nanos component (-1)")))
		})
	})

	Describe("func EqualTo()", func() {
		DescribeTable(
			"it returns true if the amounts have the same magnitude",
			func(a, b *money.Money, expect int) {
				Expect(EqualTo(a, b)).To(Equal(expect == 0))
				Expect(EqualTo(b, a)).To(Equal(expect == 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				EqualTo(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func IdenticalTo()", func() {
		DescribeTable(
			"it returns true if the amounts have the same currency and magnitude",
			func(a, b *money.Money, expect int) {
				Expect(IdenticalTo(a, b)).To(Equal(expect == 0))
				Expect(IdenticalTo(b, a)).To(Equal(expect == 0))
			},
			vectors...,
		)

		It("returns false if the amounts have a different currency", func() {
			a := &money.Money{CurrencyCode: "XYZ"}
			b := &money.Money{CurrencyCode: "ABC"}
			Expect(IdenticalTo(a, b)).To(BeFalse())
		})
	})

	Describe("func LessThan()", func() {
		DescribeTable(
			"it returns true if a < b",
			func(a, b *money.Money, expect int) {
				Expect(LessThan(a, b)).To(Equal(expect < 0))
				Expect(LessThan(b, a)).To(Equal(expect > 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				LessThan(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func LessThanOrEqualTo()", func() {
		DescribeTable(
			"it returns true if a <= b",
			func(a, b *money.Money, expect int) {
				Expect(LessThanOrEqualTo(a, b)).To(Equal(expect <= 0))
				Expect(LessThanOrEqualTo(b, a)).To(Equal(expect >= 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				LessThanOrEqualTo(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func GreaterThan()", func() {
		DescribeTable(
			"it returns true if a > b",
			func(a, b *money.Money, expect int) {
				Expect(GreaterThan(a, b)).To(Equal(expect > 0))
				Expect(GreaterThan(b, a)).To(Equal(expect < 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				GreaterThan(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func GreaterThanOrEqualTo()", func() {
		DescribeTable(
			"it returns true if a >= b",
			func(a, b *money.Money, expect int) {
				Expect(GreaterThanOrEqualTo(a, b)).To(Equal(expect >= 0))
				Expect(GreaterThanOrEqualTo(b, a)).To(Equal(expect <= 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := &money.Money{CurrencyCode: "XYZ"}
				b := &money.Money{CurrencyCode: "ABC"}
				GreaterThanOrEqualTo(a, b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})
})

var _ = Describe("func LexicallyLessThan()", func() {
	It("allows for lexical sorting of amounts by currency, then value", func() {
		shuffled := []*money.Money{
			{CurrencyCode: "XXA", Units: 1},
			{CurrencyCode: "XXB", Units: -1},
			{CurrencyCode: "XXB", Units: 1},
			{CurrencyCode: "XXB", Units: 0},
			{CurrencyCode: "XXA", Units: 0},
			{CurrencyCode: "XXA", Units: -1},
		}

		sorted := []*money.Money{
			{CurrencyCode: "XXA", Units: -1},
			{CurrencyCode: "XXA", Units: 0},
			{CurrencyCode: "XXA", Units: 1},
			{CurrencyCode: "XXB", Units: -1},
			{CurrencyCode: "XXB", Units: 0},
			{CurrencyCode: "XXB", Units: 1},
		}

		sort.Slice(
			shuffled,
			func(i, j int) bool {
				return LexicallyLessThan(shuffled[i], shuffled[j])
			},
		)

		Expect(shuffled).To(Equal(sorted))
	})
})

var _ = Describe("func Min()", func() {
	It("returns the smallest of the given amounts", func() {
		Expect(
			Min(
				&money.Money{CurrencyCode: "XYZ", Units: 2},
				&money.Money{CurrencyCode: "XYZ", Units: 1},
				&money.Money{CurrencyCode: "XYZ", Units: 3},
			),
		).To(Equal(
			&money.Money{CurrencyCode: "XYZ", Units: 1},
		))
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Min()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Min(
				&money.Money{CurrencyCode: "XYZ"},
				&money.Money{CurrencyCode: "ABC"},
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (ABC vs XYZ)"))
	})
})

var _ = Describe("func Max()", func() {
	It("returns the largest of the given amounts", func() {
		Expect(
			Max(
				&money.Money{CurrencyCode: "XYZ", Units: 2},
				&money.Money{CurrencyCode: "XYZ", Units: 1},
				&money.Money{CurrencyCode: "XYZ", Units: 3},
			),
		).To(Equal(
			&money.Money{CurrencyCode: "XYZ", Units: 3},
		))
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Max()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Max(
				&money.Money{CurrencyCode: "XYZ"},
				&money.Money{CurrencyCode: "ABC"},
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (ABC vs XYZ)"))
	})
})
