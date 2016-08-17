package mm

import "reflect"

// OriginalMethod ...
type OriginalMethod struct {
	name string
	m    reflect.Value
}

// OriginalMethodNew ...
func OriginalMethodNew(methodName string, method interface{}) (o *OriginalMethod, ok bool) {
	if !IsMethod(method) {
		ok = false
		return
	}
	v := OriginalMethod{
		name: methodName,
		m:    reflect.ValueOf(method),
	}
	o = &v
	ok = true
	return
}

// Apply ...
func (o *OriginalMethod) Apply(input []reflect.Value) (output []reflect.Value) {

	fullKey := getFullKey(o.m.Interface())
	respQueue := GetLocalRespQueue(fullKey)
	element, ok := respQueue.Shift()

	if !ok {
		// fall back to original method if resp queue is empty
		output = o.m.Call(input)

		if currentLogger != nil {
			outputElements := ValueSliceToInterfaceSlice(output)
			currentLogger(o.name, false, outputElements)
		}

	} else if fakeFunc, yes := isTempImpl(element); yes {

		f := reflect.ValueOf(fakeFunc.impl)

		// call the fakefunc and reply
		output = f.Call(input)

		if currentLogger != nil {
			outputElements := ValueSliceToInterfaceSlice(output)
			currentLogger(o.name, true, outputElements)
		}

	} else {

		// use first resp in the queue
		t := o.m.Type()
		output = make([]reflect.Value, len(element))
		for i, v := range element {
			if v == nil {
				output[i] = reflect.Zero(t.Out(i))
			} else {
				output[i] = reflect.ValueOf(v).Convert(t.Out(i))
			}
		}

		if currentLogger != nil {
			outputElements := ValueSliceToInterfaceSlice(output)
			currentLogger(o.name, true, outputElements)
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
	f = reflect.MakeFunc(o.m.Type(), o.Apply)
	return
}
