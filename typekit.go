package typekit

import (
	"errors"
	"fmt"
	"reflect"
)

type registerable interface {
	// indicates whether or not the struct is initialized
	isRegistered() bool

	// deinits the struct. Should always panic if Initialized()
	// is true.
	Unregister()
}

// this maps the static types to the instanceMap that are registered
var instanceMap = make(map[string]registerable)

func lookupInstance[T any]() (b instance[T], key string, exists bool, beantype reflect.Type) {
	val := *new(T)
	beantype = reflect.TypeOf(val)
	key = beantype.PkgPath() + beantype.Name()
	b, exists = instanceMap[key].(instance[T])
	return b, key, exists, beantype
}

// each instance holds a pointer to a struct of
// the type that it corresponds to.
type instance[T any] struct {
	// the pointer to the instance's struct
	val *T
	// true or false, whether the instance is initialized or not
	initialized bool
}

func (b instance[T]) isRegistered() bool {
	return b.initialized
}

func (b instance[T]) Unregister() {
	if !b.isRegistered() {
		panic(fmt.Errorf("cannot de-initialize an instance that hasn't been de-initialized"))
	}
	b.initialized = false
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
		if !instance.isRegistered() {
			panic(fmt.Errorf("cannot get instance: instance for type %s has been unregistered", instancetype.Name()))
		}
		// if exists, return the value
		return instance.val
	} else {
		panic(fmt.Errorf("cannot get instance: type %s has no registered instance. Ensure that the package defining the type also registers it as an instance with typekit.Register()", instancetype))
	}
}

/*
Registers an instance of the specified type. The
provided parameters can be either:

 1. An instance to store globally (typekit was designed for
    use with structs, but pointers can also work here).

 2. The result of a call to a function that returns an
    instance (see #1) and an error. In this case, if the
    error is non-nil, Register() will panic.

Example for case 1:

	// this simply assigns the instance to the registered instance
	var myPackageLevelVar = typekit.Register(MyCoolType{})

Example for case 2:

	// here, the error returned by SomeFunc() will be checked.
	// If non-nil, Register() will panic.
	var myPackageLevelVar = typekit.Register(SomeFunc())
*/
func Register[T any](val T, err ...error) *T {
	// check if an instance of the specified type exists
	instance_, key, exists, instancetype := lookupInstance[T]()
	// if the user called a function from within the
	// function parentheses, the error must be checked
	if len(err) != 0 && err[0] != nil {

		panic(fmt.Errorf("error registering instance for type %s - %s", instancetype.Name(), err[0]))
	}
	if !exists {
		// instantiate a pointer to a variable of the specified type
		addr := new(T)
		// TODO: try to replace the following 2 lines with a simple `*addr = &val`
		*addr = val
		// set the instance in the instance map for the provided type
		instanceMap[key] = instance[T]{
			val:         addr,
			initialized: true,
		}
		return addr
	} else if instance_.isRegistered() {
		panic(fmt.Errorf("cannot register instance: instance for %s has already been registered", instancetype.Name()))
	} else {
		// this sets the struct but does not override the pointer
		*instance_.val = val
		return instance_.val
	}
}

/*
Unregisters an instance of the specified type. Should be
used followed with Inject(), in order to properly set the
dependencies after unregistering them
*/
func Unregister[T any]() {
	// check if an instance of the specified type exists
	instance, _, exists, instancetype := lookupInstance[T]()
	if !exists {
		panic(fmt.Errorf("cannot unregister instance: instance for %s has not been registered", instancetype.Name()))
	}
	instance.Unregister()
}

func Inject[T any](val T) {
	// check if an instance of the specified type exists
	bn, _, exists, _ := lookupInstance[T]()
	if !exists {
		panic(fmt.Errorf("cannot mock an instance that has not been registered"))
	} else if bn.isRegistered() {
		panic(errors.New("must unregister an instance in order to mock it"))
	}
	*bn.val = val
	bn.initialized = true
}
