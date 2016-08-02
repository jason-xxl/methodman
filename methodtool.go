package methodman

import (
	"fmt"
	"reflect"
)

// MethodUniqueID is string representation of a method's identifier (memory address)
// in current running instance. It's not guaranteed to be the same after process restarts.
// For being easy to understand it use fmt pkg and a string form instead of
// unsafe.Pointer
type MethodUniqueID string

// GetMethodUniqueID ...
func GetMethodUniqueID(method interface{}) (id MethodUniqueID) {
	v := reflect.ValueOf(method)

	// indirect if it's a pointer
	for {
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		} else {
			break
		}
	}

	// check if it's really a method
	if v.Kind() != reflect.Func {
		panic("GetMethodUniqueID: input param is not a method.")
	}

	// use fmt to get address of the value
	id = MethodUniqueID(fmt.Sprintf("%p", v.Interface()))
	return
}

// IsMethodPointer ...
func IsMethodPointer(methodPointer interface{}) (yes bool) {
	v := reflect.ValueOf(methodPointer)
	yes = v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Func
	return
}

// IsMethod ...
func IsMethod(method interface{}) (yes bool) {
	v := reflect.ValueOf(method)
	yes = v.Kind() == reflect.Func
	return
}

// GetMethodType ...
func GetMethodType(method interface{}) (typeMethod reflect.Type) {
	t := reflect.TypeOf(method)

	// indirect if it's a pointer
	for {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		} else {
			break
		}
	}

	// check if it's really a method
	if t.Kind() != reflect.Func {
		panic("GetMethodUniqueID: input param is not a method.")
	}

	return
}
