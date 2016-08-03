package mm

import "reflect"

// EnableMock upgrades a method to allow per-goroutine mocking.
// Note that,
// 1) it will replace your method with method helper, so do it at beginning of all your
// unittests, not in the middle, to avoid race condition.
// 2) if default queue length (200) is not enough for you, you can enlarge by giving the
// queueLength. It panic if you gave invalid value.
// 3) it's expected to be called at ONLY init phase of tests, so no need to protect with
// Mutex
func EnableMock(method interface{}, name string) {

	if !IsMethodPointer(method) {
		panic("methodman.EnableMock: method is not a method pointer.")
	}

	// assuming its a manager method, check if it is
	methodKey := GetMethodUniqueID(method)
	if _, ok := managerMap[methodKey]; ok {
		// do nothing if methodKey exists
		return
	}

	// if it's an original method, assign a manager for it
	v := reflect.ValueOf(method)
	realMethod := v.Elem().Interface()
	manager := ManagerNew(name, realMethod)

	// create a manager method for the original method
	f := manager.Method.MakeFunc()

	// store the manager method
	methodKey = GetMethodUniqueID(f.Interface())
	managerMap[methodKey] = manager

	// switch the pointer of original method to the manager method
	v.Elem().Set(f)
	return
}
