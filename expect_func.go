package mm

import "reflect"

// tempImpl is used to tag the fake response as a temp implementation, not a return value
type tempImpl struct {
	impl interface{}
}

// ExpectFunc adds a temp implementation of the original method and consume once
// It would be helpful for some special case like simulating a timeout.
func ExpectFunc(method interface{}, fakeFunc interface{}) {

	manager, respQueue := getQueueFromMethod(method)

	// signature check
	typeOriginalMethod := manager.Method.m.Type()
	typeTempImpl := reflect.TypeOf(fakeFunc)

	if typeOriginalMethod != typeTempImpl {
		panic("methodman.ExpectFunc: fakeFunc's signature doesn't match original method.")
	}

	// push to queue with a special type
	ok := respQueue.Push([]interface{}{
		tempImpl{
			impl: fakeFunc,
		},
	})

	if !ok {
		panic("methodman.ExpectFunc: " + errorInfoQueueFull)
	}
	return
}
