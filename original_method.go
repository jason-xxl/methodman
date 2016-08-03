package mm

import "reflect"

// OriginalMethod ...
type OriginalMethod reflect.Value

// OriginalMethodNew ...
func OriginalMethodNew(method interface{}) (o *OriginalMethod, ok bool) {
	if !IsMethod(method) {
		ok = false
		return
	}
	v := OriginalMethod(reflect.ValueOf(method))
	o = &v
	ok = true
	return
}

// Apply ...
func (o *OriginalMethod) Apply(input []reflect.Value) (output []reflect.Value) {

	fullKey := getFullKey(reflect.Value(*o).Interface())
	respQueue := GetLocalRespQueue(fullKey)
	element, ok := respQueue.Shift()

	if !ok {
		// fall back to original method if resp queue is empty
		output = reflect.Value(*o).Call(input)

	} else if fakeFunc, yes := isTempImpl(element); yes {

		f := reflect.ValueOf(fakeFunc.impl)

		// call the fakefunc and reply
		output = f.Call(input)

	} else {
		// use first resp in the queue
		t := reflect.Value(*o).Type()
		output = make([]reflect.Value, len(element))
		for i, v := range element {
			if v == nil {
				output[i] = reflect.Zero(t.Out(i))
			} else {
				output[i] = reflect.ValueOf(v).Convert(t.Out(i))
			}
		}
	}

	return
}

func isTempImpl(element []interface{}) (f tempImpl, yes bool) {
	if len(element) != 1 {
		return
	}
	f, yes = element[0].(tempImpl)
	if f.impl == nil {
		yes = false
	}
	return
}

// MakeFunc ...
func (o *OriginalMethod) MakeFunc() (f reflect.Value) {
	f = reflect.MakeFunc(reflect.Value(*o).Type(), o.Apply)
	return
}
