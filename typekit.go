package typekit

import (
	"errors"
	"fmt"
	"reflect"
)

type initializeableBean interface {
	// indicates whether or not the struct is initialized
	isRegistered() bool

	// deinits the struct. Should always panic if Initialized()
	// is true.
	Unregister()
}

// this maps the static types to the beansMap that are registered
var beansMap = make(map[string]initializeableBean)

func lookupBeanInMap[T any](val ...T) (b bean[T], key string, exists bool, beantype reflect.Type) {
	if len(val) == 0 {
		val = append(val, *new(T))
	}
	beantype = reflect.TypeOf(val)
	key = beantype.PkgPath() + beantype.Name()
	b, exists = beansMap[key].(bean[T])
	return b, key, exists, beantype
}

// each bean holds the struct itself, but
type bean[T any] struct {
	// the pointer to the bean's struct
	val *T
	// true or false, whether the bean is initialized or not
	initialized bool
}

func (b bean[T]) isRegistered() bool {
	return b.initialized
}

func (b bean[T]) Unregister() {
	if !b.isRegistered() {
		panic(fmt.Errorf("cannot de-initialize a bean that hasn't been de-initialized"))
	}
	b.initialized = false
}

func Get[T any]() *T {
	bn, _, exists, beantype := lookupBeanInMap[T]()
	if exists {
		if !bn.isRegistered() {
			panic(fmt.Errorf("cannot get bean: bean %s has not been initialized or has been de-initialized", beantype.Name()))
		}
		// if exists, return the value
		return bn.val
	} else {
		panic(fmt.Errorf("no bean has been registered with type %s. Ensure that the package defining the type also registers it as a bean with typekit.Register()", beantype))
	}
}

/*
 */
func Register[T any](val T) *T {
	// check if a bean of the specified type exists
	bn, key, exists, beantype := lookupBeanInMap[T]()
	if !exists {
		// instantiate a pointer to a variable of the specified type
		addr := new(T)
		// TODO: try to replace the following 2 lines with a simple `*addr = &val`
		*addr = val
		// set the bean in the bean map for the provided type
		beansMap[key] = bean[T]{
			val:         addr,
			initialized: true,
		}
		return addr
	} else if bn.isRegistered() {
		panic(fmt.Errorf("cannot register bean: bean %s has already been registered", beantype.Name()))
	} else {
		// this sets the struct but does not override the pointer
		*bn.val = val
		return bn.val
	}
}

func Unregister[T any]() {
	// check if a bean of the specified type exists
	bn, _, exists, beantype := lookupBeanInMap[T]()
	if !exists {
		panic(fmt.Errorf("cannot unregister bean: bean %s has not been registered", beantype.Name()))
	}
	bn.Unregister()
}

func Mock[T any](val T) {
	// check if a bean of the specified type exists
	bn, _, exists, _ := lookupBeanInMap(val)
	if !exists {
		panic(fmt.Errorf("cannot mock a bean that has not been registered"))
	} else if bn.isRegistered() {
		panic(errors.New("must unregister a bean in order to mock it"))
	}
	*bn.val = val
}
