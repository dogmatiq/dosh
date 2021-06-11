package dosh_test

import (
	"math"

	. "github.com/dogmatiq/dosh"
	. "github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
)

var _ = Describe("type Amount (protocol buffers marshaling)", func() {
	Describe("func MarshalProto()", func() {
		It("returns the protocol buffers representation of the amount", func() {
			a := MustParse("XYZ", "10.123")

			pb, err := a.MarshalProto()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(pb).To(EqualX(&money.Money{
				CurrencyCode: "XYZ",
				Units:        10,
				Nanos:        123000000,
			}))
		})

		DescribeTable(
			"it returns an error if the amount can not be marshaled",
			func(a Amount, expect string) {
				_, err := a.MarshalProto()
				Expect(err).To(MatchError(expect))
			},
			Entry(
				"integer component of the magnitude overflows an int64",
				Int("XYZ", math.MaxInt64).Add(Unit("XYZ")),
				"cannot marshal amount to protocol buffers representation: magnitude's integer component overflows int64",
			),
			Entry(
				"fractional component of the magnitude requires more precision than available",
				MustParse("XYZ", "0.0123456789"),
				"cannot marshal amount to protocol buffers representation: magnitude's fractional component has too many decimal places",
			),
		)
	})

	Describe("func UnmarshalProto()", func() {
		It("unmarshals an amount from its protocol buffers representation", func() {
			var a Amount

			err := a.UnmarshalProto(&money.Money{
				CurrencyCode: "XYZ",
				Units:        10,
				Nanos:        123000000,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(a.CurrencyCode()).To(Equal("XYZ"))

			m := decimal.RequireFromString("10.123")
			Expect(a.Magnitude().Equal(m))
		})

		DescribeTable(
			"it returns an error if the protocol buffers message is invalid",
			func(pb *money.Money, expect string) {
				var a Amount
				err := a.UnmarshalProto(pb)
				Expect(err).To(MatchError(expect))
			},
			Entry(
				"empty currency",
				&money.Money{},
				"cannot unmarshal amount from protocol buffers representation: currency code is empty, codes must consist only of 3 or more uppercase ASCII letters",
			),
			Entry(
				"invalid currency",
				&money.Money{CurrencyCode: "X"},
				"cannot unmarshal amount from protocol buffers representation: currency code (X) is invalid, codes must consist only of 3 or more uppercase ASCII letters",
			),
			Entry(
				"units positive, nanos negative",
				&money.Money{CurrencyCode: "XYZ", Units: +1, Nanos: -1},
				"cannot unmarshal amount from protocol buffers representation: units and nanos components must have the same sign",
			),
			Entry(
				"units negative, nanos positive",
				&money.Money{CurrencyCode: "XYZ", Units: -1, Nanos: +1},
				"cannot unmarshal amount from protocol buffers representation: units and nanos components must have the same sign",
			),
		)
	})
})
