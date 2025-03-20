package internal

import (
	"fmt"
	"reflect"
)

type initializable interface {
	// indicates whether or not the struct is initialized
	isInitialized() bool

	// initializes the struct
	initialize()

	// de-inits the struct
	deinit()

	// un-mocks the struct
	unmock()
}

// this maps the static types to the instanceMap that are registered
var instanceMap = make(map[string]initializable)

func lookupInstance[T any]() (b *instance[T], key string, exists bool, beantype reflect.Type) {
	val := *new(T)
	beantype = reflect.TypeOf(val)
	key = beantype.PkgPath() + beantype.Name()
	b, exists = instanceMap[key].(*instance[T])
	return b, key, exists, beantype
}

// each instance holds a pointer to a struct of
// the type that it corresponds to.
type instance[T any] struct {
	// the pointer to the instance's struct
	val *T
	// true or false, whether the instance is initialized or not
	initialized bool
	// the function that marks the init function for that instance type
	initFn func() (T, error)

	// true or false, whether the instance is mocked or not
	mocked bool
	// this function is to be run when initialize() is called if mocked is true
	mockInit func() (T, error)
}

func (b *instance[T]) isInitialized() bool {
	return b.initialized
}

func (b *instance[T]) initialize() {
	var err error
	if b.mocked {
		*b.val, err = b.mockInit()
	} else {
		*b.val, err = b.initFn()
	}
	b.initialized = true
	if err != nil {
		panic(fmt.Errorf("could not initialize instance for type %T - %w", *b.val, err))
	}
}

func (b *instance[T]) deinit() {
	b.initialized = false
}

func (b *instance[T]) mock(initFn func() (T, error)) {

}

func (b *instance[T]) unmock() {
	b.mocked = false
	b.mockInit = nil
}

/*
Resolves the instance of the specified type. If the instance
has not been registered, or if it has been unregistered,
this panics.

Example:

	// in this example, the instance of SomeType is made accessible
	// from everywhere else in the package, through the
	// somePackageLevelVar pointer.
	var somePackageLevelVar = typekit.Resolve[somepackage.SomeType]()
*/
func Resolve[T any]() *T {
	instance, _, exists, instancetype := lookupInstance[T]()
	if exists {
		if !instance.isInitialized() {
			instance.initialize()
		}
		// if exists, return the value
		return instance.val
	} else {
		panic(fmt.Errorf("cannot get instance: type %s has no registered instance. Ensure that the package defining the type also registers it as an instance with typekit.Register()", instancetype))
	}
}

/*
Registers an instance of the specified type. The
provided function will be executed when needed, based
on the implicit dependency tree between types registered
using typekit.

Example:

	package foo

	import "github.com/the-zucc/somepackage"

	type Foo struct {}

	var myPackageLevelVar = typekit.Register(func() (Foo, error) {
		return Foo{}, nil
	})

	// ...in another package:

	package bar

	import "github.com/the-zucc/foo"

	type Bar struct {
		field *Foo
	}
	var foobar = typekit.Register(func() (Bar, error){
		return Bar{
			field: typekit.Resolve[foo.Foo]() // this will be executed in "lazy" mode
		}
	})
*/
func Register[T any](initFn func() (T, error)) *T {
	// check if an instance of the specified type exists
	inst, key, exists, instancetype := lookupInstance[T]()

	// if the user called a function from within the
	// function parentheses, the error must be checked
	if !exists {
		// instantiate a pointer to a variable of the specified type
		addr := new(T)
		// set the instance in the instance map for the provided type
		instanceMap[key] = &instance[T]{
			val:         addr,
			initialized: true,
			initFn:      initFn,
		}
		return addr
	} else if inst.isInitialized() {
		panic(fmt.Errorf("cannot register instance: instance for %s has already been registered", instancetype.Name()))
	} else {
		// this sets the struct but does not override the pointer
		return inst.val
	}
}

/*
Used to inject dependencies after they have been registered.
Should be used for creating mocks and other testing things.
*/
func Mock[T any](f func() (T, error)) {
	// check if an instance of the specified type exists
	bn, _, exists, _ := lookupInstance[T]()
	if !exists {
		panic(fmt.Errorf("cannot mock an instance that has not been registered"))
	}
	bn.mockInit = f
}

/*
Used to refresh the dependency tree. Essentially, it marks all
instances as needing a re-initialization. In the proccess, instances
which have been mocked using mock.Mock() will be initialized using
their mock constructors.
*/
func RefreshTree() {
	for _, val := range instanceMap {
		// this will have the next call to Resolve() with that type re-run
		// initialization function or the mock, whichever case applies
		val.deinit()
	}
}

func ResetMocks() {
	for _, val := range instanceMap {
		val.unmock()
	}
}

func ResetMock[T any]() {
	a, _, exists, _ := lookupInstance[T]()
	if !exists {
		panic(fmt.Errorf("could not reset mock for type %T - no instance of that type was registered", a.val))
	}
}
