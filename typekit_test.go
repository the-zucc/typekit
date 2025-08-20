package typekit_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/the-zucc/typekit"
)

func TestTypeKit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TypeKit Suite")
}

type Foo interface {
}
type Bar struct{}

func (b Bar) Foo() {}

var _ = Describe("typekit", func() {
	It("should work", func() {
		_ = typekit.Register(func() (Foo, error) {
			return Bar{}, nil
		})

		Expect(typekit.Resolve[Foo]()).To(BeAssignableToTypeOf(Foo(nil)))
		Expect(true).To(BeTrue())
	})
})
