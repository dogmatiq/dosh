package dosh_test

import (
	"strings"

	. "github.com/dogmatiq/dosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

var _ = Describe("type Amount (binary marshaling)", func() {
	Describe("func MarshalBinary() and UnmarshalBinary()", func() {
		It("marshals and unmarshals an amount", func() {
			a := FromString("XYZ", "10.123")

			data, err := a.MarshalBinary()
			Expect(err).ShouldNot(HaveOccurred())

			var b Amount
			err = b.UnmarshalBinary(data)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(a.EqualTo(b)).To(BeTrue())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("returns an error if the currency code is longer than 255 bytes", func() {
			a := Zero(strings.Repeat("X", 256))
			_, err := a.MarshalBinary()
			Expect(err).To(MatchError("cannot marshal amount to binary representation: currency code is 256 characters long, maximum is 255"))
		})
	})

	// binaryZero is the binary representation of a decimal with a value of 0.
	//
	// This is built inline (as opposed to in BeforeEach) because it is used
	// within a call to Entry().
	var binaryZero string

	{
		data, err := decimal.Decimal{}.MarshalBinary()
		if err != nil {
			panic(err)
		}

		binaryZero = string(data)
	}

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
				"\x00"+binaryZero,
				"cannot unmarshal amount from binary representation: currency code is empty, codes must consist only of 3 or more uppercase ASCII letters",
			),
			Entry(
				"invalid currency",
				"\x01X"+binaryZero,
				"cannot unmarshal amount from binary representation: currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters",
			),
			Entry(
				"empty magnitude",
				"\x03USD",
				"cannot unmarshal amount from binary representation: error decoding binary []: expected at least 5 bytes, got 0",
			),
			Entry(
				"invalid magnitude",
				"\x03USD<invalid>",
				"cannot unmarshal amount from binary representation: error decoding binary [60 105 110 118 97 108 105 100 62]: Int.GobDecode: encoding version 48 not supported",
			),
		)
	})
})
