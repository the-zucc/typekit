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

func RefreshTree() {
	internal.RefreshTree()
}

func ResetMocks() {
	internal.ResetMocks()
}

func ResetMock[T any]() {
	internal.ResetMock[T]()
}
