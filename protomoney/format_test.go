package protomoney_test

import (
	"fmt"

	. "github.com/dogmatiq/dosh/protomoney"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/type/money"
)

func ExampleFmt() {
	m := &money.Money{
		CurrencyCode: "XYZ",
		Units:        10,
		Nanos:        129000000,
	}

	fmt.Printf("%.2f\n", Fmt(m))
	// Output: XYZ 10.13
}

var _ = Describe("func Fmt()", func() {
	It("returns a formatted representation of the amount", func() {
		m := &money.Money{
			CurrencyCode: "XYZ",
			Units:        10,
			Nanos:        129000000,
		}
		s := fmt.Sprintf("%0.2f", Fmt(m))
		Expect(s).To(Equal("XYZ 10.13"))
	})

	It("returns a descriptive string if used with an unsupported verb", func() {
		m := &money.Money{
			CurrencyCode: "XYZ",
			Units:        10,
			Nanos:        129000000,
		}
		s := fmt.Sprintf("%d", Fmt(m))
		Expect(s).To(Equal("%!d(*money.Money=" + m.String() + ")"))
	})
})
