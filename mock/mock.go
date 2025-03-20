package mock

import "github.com/the-zucc/typekit/internal"

func Inject[T any](val T) {
	internal.Inject(val)
}
