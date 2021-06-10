package dosh_test

import (
	"sort"

	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Amount (comparison functions)", func() {
	vectors := []TableEntry{
		Entry("equal, zero", "0", "0", 0),
		Entry("equal, positive integer", "+1", "+1", 0),
		Entry("equal, negative integer", "-1", "-1", 0),
		Entry("equal, positive decimal, no units", "+0.123", "+0.123", 0),
		Entry("equal, negative decimal, no units", "-0.123", "-0.123", 0),
		Entry("equal, positive decimal", "+1.123", "+1.123", 0),
		Entry("equal, negative decimal", "-1.123", "-1.123", 0),

		Entry("less, zero", "0", "+1", -1),
		Entry("less, positive integer", "+1", "+2", -1),
		Entry("less, negative integer", "-2", "-1", -1),
		Entry("less, positive decimal, no units", "+0.123", "+0.456", -1),
		Entry("less, negative decimal, no units", "-0.456", "-0.123", -1),
		Entry("less, positive decimal", "+1.123", "+1.456", -1),
		Entry("less, negative decimal", "-1.456", "-1.123", -1),
	}

	Describe("func IsZero()", func() {
		DescribeTable(
			"returns true if the amount has a magnitude of zero",
			func(dec string, expect bool) {
				Expect(MustParse("XYZ", dec).IsZero()).To(Equal(expect))
			},
			Entry("zero", "0", true),
			Entry("positive", "1.23", false),
			Entry("negative", "-1.23", false),
		)
	})

	Describe("func IsPositive()", func() {
		DescribeTable(
			"returns true if the amount has a positive magnitude",
			func(dec string, expect bool) {
				Expect(MustParse("XYZ", dec).IsPositive()).To(Equal(expect))
			},
			Entry("zero", "0", false),
			Entry("positive", "1.23", true),
			Entry("negative", "-1.23", false),
		)
	})

	Describe("func IsNegative()", func() {
		DescribeTable(
			"returns true if the amount has a negative magnitude",
			func(dec string, expect bool) {
				Expect(MustParse("XYZ", dec).IsNegative()).To(Equal(expect))
			},
			Entry("zero", "0", false),
			Entry("positive", "1.23", false),
			Entry("negative", "-1.23", true),
		)
	})

	Describe("func Cmp()", func() {
		DescribeTable(
			"it returns a C-style comparison result",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				x := l.Cmp(r)
				y := r.Cmp(l)

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
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.Cmp(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func EqualTo()", func() {
		DescribeTable(
			"it returns true if the amounts have the same magnitude",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.EqualTo(r)).To(Equal(expect == 0))
				Expect(r.EqualTo(l)).To(Equal(expect == 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.EqualTo(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func IdenticalTo()", func() {
		DescribeTable(
			"it returns true if the amounts have the same currency and magnitude",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.IdenticalTo(r)).To(Equal(expect == 0))
				Expect(r.IdenticalTo(l)).To(Equal(expect == 0))
			},
			vectors...,
		)

		It("returns false if the amounts have a different currency", func() {
			a := MustParse("XYZ", "1")
			b := MustParse("ABC", "1")
			Expect(a.IdenticalTo(b)).To(BeFalse())
		})
	})

	Describe("func LessThan()", func() {
		DescribeTable(
			"it returns true if a < b",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.LessThan(r)).To(Equal(expect < 0))
				Expect(r.LessThan(l)).To(Equal(expect > 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.LessThan(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func LessThanOrEqualTo()", func() {
		DescribeTable(
			"it returns true if a <= b",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.LessThanOrEqualTo(r)).To(Equal(expect <= 0))
				Expect(r.LessThanOrEqualTo(l)).To(Equal(expect >= 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.LessThanOrEqualTo(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func GreaterThan()", func() {
		DescribeTable(
			"it returns true if a > b",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.GreaterThan(r)).To(Equal(expect > 0))
				Expect(r.GreaterThan(l)).To(Equal(expect < 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.GreaterThan(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func GreaterThanOrEqualTo()", func() {
		DescribeTable(
			"it returns true if a >= b",
			func(a, b string, expect int) {
				l := MustParse("XYZ", a)
				r := MustParse("XYZ", b)
				Expect(l.GreaterThanOrEqualTo(r)).To(Equal(expect >= 0))
				Expect(r.GreaterThanOrEqualTo(l)).To(Equal(expect <= 0))
			},
			vectors...,
		)

		It("panics if the amounts do not have the same currency", func() {
			Expect(func() {
				a := MustParse("XYZ", "1")
				b := MustParse("ABC", "1")
				a.GreaterThanOrEqualTo(b)
			}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
		})
	})

	Describe("func LexicallyLessThan()", func() {
		It("allows for lexical sorting of amounts by currency, then value", func() {
			shuffled := []Amount{
				MustParse("XXA", "1"),
				MustParse("XXB", "-1"),
				MustParse("XXB", "1"),
				MustParse("XXB", "0"),
				MustParse("XXA", "0"),
				MustParse("XXA", "-1"),
			}

			sorted := []Amount{
				MustParse("XXA", "-1"),
				MustParse("XXA", "0"),
				MustParse("XXA", "1"),
				MustParse("XXB", "-1"),
				MustParse("XXB", "0"),
				MustParse("XXB", "1"),
			}

			sort.Slice(
				shuffled,
				func(i, j int) bool {
					return shuffled[i].LexicallyLessThan(shuffled[j])
				},
			)

			Expect(shuffled).To(Equal(sorted))
		})
	})
})

var _ = Describe("func Min()", func() {
	It("returns the smallest of the given amounts", func() {
		Expect(
			Min(
				MustParse("XYZ", "2"),
				MustParse("XYZ", "1"),
				MustParse("XYZ", "3"),
			).EqualTo(
				MustParse("XYZ", "1"),
			),
		).To(BeTrue())
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Min()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Min(
				MustParse("XYZ", "1"),
				MustParse("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})

var _ = Describe("func Max()", func() {
	It("returns the largest of the given amounts", func() {
		Expect(
			Max(
				MustParse("XYZ", "2"),
				MustParse("XYZ", "1"),
				MustParse("XYZ", "3"),
			).EqualTo(
				MustParse("XYZ", "3"),
			),
		).To(BeTrue())
	})

	It("panics if no amounts are provided", func() {
		Expect(func() {
			Max()
		}).To(PanicWith("at least one amount must be provided"))
	})

	It("panics if the amounts do not have the same currency", func() {
		Expect(func() {
			Max(
				MustParse("XYZ", "1"),
				MustParse("ABC", "1"),
			)
		}).To(PanicWith("can not operate on amounts in differing currencies (XYZ vs ABC)"))
	})
})
