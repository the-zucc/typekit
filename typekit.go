package typekit

import (
	"github.com/the-zucc/typekit/internal"
)

/*
Registers an instance of the specified type. The
provided function will be executed when needed, based
on the implicit dependency tree between types registered
using typekit.

Example:

	package foo

	import (
		"github.com/the-zucc/typekit"
		"github.com/the-zucc/somepackage"
	)

	type Foo struct {}

	var myPackageLevelVar = typekit.RegisterWith(func() (Foo, error) {
		return Foo{}, nil
	})

	// ...in another package:

	package bar

	import (
		"github.com/the-zucc/foo"
		"github.com/the-zucc/typekit"
	)

	type Bar struct {
		field *Foo
	}
	var foobar = typekit.RegisterWith(func() (Bar, error){
		return Bar{
			field: typekit.Resolve[foo.Foo]() // this will be executed in "lazy" mode
		}
	})
*/
func Register[T any](initFn func() (T, error)) *T {
	return internal.Register(initFn)
}

/*
Resolves the instance of the specified type. If the instance
has not been registered, or if it has been unregistered,
this panics.

Example:

	// in this example, the instance of SomeType is made accessible
	// from everywhere else in the package, through the
	// somePackageLevelVar pointer.
	var someVar = typekit.Resolve[somepackage.SomeType]()
*/
func Resolve[T any]() *T {
	return internal.Resolve[T]()
}
