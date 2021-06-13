package protomoney_test

import (
	. "github.com/dogmatiq/dosh/protomoney"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/type/money"
)

var _ = Describe("func Validate()", func() {
	DescribeTable(
		"it returns nil if the amount is valid",
		func(m *money.Money) {
			err := Validate(m)
			Expect(err).ShouldNot(HaveOccurred())
		},
		Entry("zero", &money.Money{CurrencyCode: "XYZ"}),
		Entry("positive units", &money.Money{CurrencyCode: "XYZ", Units: 1}),
		Entry("positive nanos", &money.Money{CurrencyCode: "XYZ", Nanos: 1}),
		Entry("positive units & nanos", &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 1}),
		Entry("negative units", &money.Money{CurrencyCode: "XYZ", Units: -1}),
		Entry("negative nanos", &money.Money{CurrencyCode: "XYZ", Nanos: -1}),
		Entry("negative units & nanos", &money.Money{CurrencyCode: "XYZ", Units: -1, Nanos: -1}),
	)

	DescribeTable(
		"it returns an error if the amount is invalid",
		func(m *money.Money, expect string) {
			err := Validate(m)
			Expect(err).To(MatchError(expect))
		},
		Entry("empty currency code", &money.Money{}, "currency code is empty, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("invalid currency code", &money.Money{CurrencyCode: "X"}, "currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("positive units, negative nanos", &money.Money{CurrencyCode: "XYZ", Units: +1, Nanos: -1}, "sign of units component (1) does not agree with sign of nanos component (-1)"),
		Entry("negative units, positive nanos", &money.Money{CurrencyCode: "XYZ", Units: -1, Nanos: +1}, "sign of units component (-1) does not agree with sign of nanos component (1)"),
	)
})

var _ = Describe("func Normalize()", func() {
	DescribeTable(
		"it normalizes the amount",
		func(m, expect *money.Money) {
			n, err := Normalize(m)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(n).To(Equal(expect))
		},
		Entry("normalized", &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 999999999}, &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 999999999}),
		Entry("denormalized, overflow", &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 1500000000}, &money.Money{CurrencyCode: "XYZ", Units: 2, Nanos: 500000000}),
		Entry("denormalized, underflow", &money.Money{CurrencyCode: "XYZ", Units: 1, Nanos: 1500000000}, &money.Money{CurrencyCode: "XYZ", Units: 2, Nanos: 500000000}),
	)

	DescribeTable(
		"it returns an error if the amount is invalid",
		func(m *money.Money, expect string) {
			_, err := Normalize(m)
			Expect(err).To(MatchError(expect))
		},
		Entry("empty currency code", &money.Money{}, "currency code is empty, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("invalid currency code", &money.Money{CurrencyCode: "X"}, "currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("positive units, negative nanos", &money.Money{CurrencyCode: "XYZ", Units: +1, Nanos: -1}, "sign of units component (1) does not agree with sign of nanos component (-1)"),
		Entry("negative units, positive nanos", &money.Money{CurrencyCode: "XYZ", Units: -1, Nanos: +1}, "sign of units component (-1) does not agree with sign of nanos component (1)"),
	)
})
