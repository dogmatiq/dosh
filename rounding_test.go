package dosh_test

import (
	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Amount (rounding methods)", func() {
	Describe("func Floor()", func() {
		DescribeTable(
			"it returns an amount with the magnitude rounded down to the nearest integer",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).Floor().EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.2", "1"),
			Entry("negative", "-1.2", "-2"),
		)
	})

	Describe("func Ceil()", func() {
		DescribeTable(
			"it returns an amount with the magnitude rounded up to the nearest integer",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).Ceil().EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.2", "2"),
			Entry("negative", "-1.2", "-1"),
		)
	})

	Describe("func Truncate()", func() {
		DescribeTable(
			"it returns an amount with the magnitude truncated to the given number of decimal places",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).Truncate(1).EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive", "1.23", "1.2"),
			Entry("negative", "-1.23", "-1.2"),
		)
	})

	Describe("func Round()", func() {
		DescribeTable(
			"it returns an amount with the magnitude rounded to the given number of decimal places",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).Round(1).EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive (down)", "1.23", "1.2"),
			Entry("positive (up)", "1.27", "1.3"),
			Entry("positive (half)", "1.25", "1.3"),
			Entry("negative (up)", "-1.23", "-1.2"),
			Entry("negative (down)", "-1.27", "-1.3"),
			Entry("negative (half)", "-1.25", "-1.3"),
		)

		DescribeTable(
			"it returns an amount with the magnitude rounded to the given number of integer places when n is negative",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).Round(-1).EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive (down)", "123", "120"),
			Entry("positive (up)", "127", "130"),
			Entry("positive (half)", "125", "130"),
			Entry("negative (up)", "-123", "-120"),
			Entry("negative (down)", "-127", "-130"),
			Entry("negative (half)", "-125", "-130"),
		)
	})

	Describe("func RoundBank()", func() {
		DescribeTable(
			"it returns an amount with the magnitude rounded to the given number of decimal places",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).RoundBank(1).EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive (down)", "1.23", "1.2"),
			Entry("positive (up)", "1.27", "1.3"),
			Entry("positive (half, even)", "1.25", "1.2"),
			Entry("positive (half, odd)", "1.15", "1.2"),
			Entry("negative (up)", "-1.23", "-1.2"),
			Entry("negative (down)", "-1.27", "-1.3"),
			Entry("negative (half, even)", "-1.25", "-1.2"),
			Entry("negative (half, odd)", "-1.15", "-1.2"),
		)

		DescribeTable(
			"it returns an amount with the magnitude rounded to the given number of integer places when n is negative",
			func(a, expect string) {
				Expect(
					FromString("XYZ", a).RoundBank(-1).EqualTo(
						FromString("XYZ", expect),
					),
				).To(BeTrue())
			},
			Entry("zero", "0", "0"),
			Entry("positive (down)", "123", "120"),
			Entry("positive (up)", "127", "130"),
			Entry("positive (half, even)", "125", "120"),
			Entry("positive (half, odd)", "115", "120"),
			Entry("negative (up)", "-123", "-120"),
			Entry("negative (down)", "-127", "-130"),
			Entry("negative (half, even)", "-125", "-120"),
			Entry("negative (half, odd)", "-115", "-120"),
		)
	})
})
