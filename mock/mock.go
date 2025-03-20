package mock

import (
	"github.com/the-zucc/typekit/internal"
)

/*
Used to inject dependencies after they have been registered.
Should be used for creating mocks and other testing things.
*/
func Mock[T any](f func() (T, error)) {
	internal.Mock(f)
}

/*
Used to refresh the dependency tree. Essentially, it marks all
instances as needing a re-initialization. In the proccess, instances
which have been mocked using mock.Mock() will be initialized using
their mock constructors.

Example:

	package foo

	import (
		"github.com/the-zucc/typekit"
		typekitmock "github.com/the-zucc/typekit/mock"
		"github.com/the-zucc/bar"
	)

	func MyTestFunction() {
		typekitmock.Mock(func() (bar.Bar, error) {
			return bar.Bar{}, nil
		})
		typekitmock.RefreshTree()
	}
*/
func RefreshTree() {
	internal.RefreshTree()
}

/*
Used to reset mocks for all instances registered using typekit.

This can be called at the beginning of a test to ensure all mocks
are reset and that the test behavior is controlled.

Example:

	package foo

	import (
		"github.com/the-zucc/typekit"
		typekitmock "github.com/the-zucc/typekit/mock"
		"github.com/the-zucc/bar"
	)

	func MyTestFunction() {
		typekitmock.ResetMocks()
		typekitmock.Mock(func() (bar.Bar, error) {
			return bar.Bar{}, nil
		})
		typekitmock.RefreshTree()
	}
*/
func ResetMocks() {
	internal.ResetMocks()
}

/*
Used to reset the mock of an instance registered using typekit.

This can be called at the beginning of a test to ensure that mock
is reset and that the test behavior is controlled for that instance.

Example:

	package foo

	import (
		"github.com/the-zucc/typekit"
		typekitmock "github.com/the-zucc/typekit/mock"
		"github.com/the-zucc/bar"
	)

	func MyTestFunction() {
		typekitmock.ResetMock[bar.Bar]()
		typekitmock.Mock(func() (bar.Bar, error) {
			return bar.Bar{}, nil
		})
		typekitmock.RefreshTree()
	}
*/
func ResetMock[T any]() {
	internal.ResetMock[T]()
}
