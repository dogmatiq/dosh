package dosh_test

import (
	"strings"

	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Amount (binary marshaling)", func() {
	Describe("func MarshalBinary() and UnmarshalBinary()", func() {
		It("marshals and unmarshals an amount", func() {
			a := MustParse("xyz", "10.123")

			data, err := a.MarshalBinary()
			Expect(err).ShouldNot(HaveOccurred())

			var b Amount
			err = b.UnmarshalBinary(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(a.Equal(b)).To(BeTrue())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("returns an error if the currency code is longer than 255 bytes", func() {
			a := Zero(strings.Repeat("X", 256))
			_, err := a.MarshalBinary()
			Expect(err).To(MatchError("cannot marshal amount to binary representation: currency code is 256 characters long, maximum is 255"))
		})
	})

	Describe("func UnmarshalBinary()", func() {
		DescribeTable(
			"it returns an error if the data is invalid",
			func(data string, expect string) {
				var a Amount
				err := a.UnmarshalBinary([]byte(data))
				Expect(err).To(MatchError(expect))
			},
			Entry(
				"empty",
				"",
				"cannot unmarshal amount from binary representation: data is empty",
			),
			Entry(
				"empty currency",
				"\x00",
				"cannot unmarshal amount from binary representation: currency component is empty",
			),
			Entry(
				"empty magnitude",
				"\x03USD",
				"cannot unmarshal amount from binary representation: data is shorter than expected",
			),
			Entry(
				"invalid magnitude",
				"\x03USD<invalid>",
				"cannot unmarshal amount from binary representation: Int.GobDecode: encoding version 48 not supported",
			),
		)
	})
})