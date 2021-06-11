package currency_test

import (
	. "github.com/dogmatiq/dosh/internal/currency"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ValidateCode()", func() {
	DescribeTable(
		"returns nil if the currency code is valid",
		func(c string) {
			err := ValidateCode(c)
			Expect(err).ShouldNot(HaveOccurred())
		},
		Entry("ISO-4217 code", "USD"),
		Entry("non-standard code", "XYZ"),
		Entry("non-standard code longer than 3 characters", "XYZXX"),
		Entry("non-standard code longer than 3 characters, without leading X", "ABCXX"),
	)

	DescribeTable(
		"returns an error if the currency code is invalid",
		func(c, expect string) {
			err := ValidateCode(c)
			Expect(err).To(MatchError(expect))
		},
		Entry("empty", "", "currency code is empty, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("too short", "X", "currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters"),
		Entry("non-letters", "XY9", "currency code (XY9) is invalid, codes must consist only of 3 or more uppercase ASCII letters"),
	)
})
