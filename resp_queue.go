package methodman

import (
	"reflect"

	"github.com/myteksi/go/commons/util/type/methodtool"
	mygls "github.com/tylerb/gls"
)

// RespQueue is a blocking queue storing expected / fake response for use
type RespQueue chan []interface{}

// RespQueueNew ...
func RespQueueNew(length int) (o RespQueue) {
	o = make(RespQueue, queueLength)
	return
}

// Push ...
func (o RespQueue) Push(element []interface{}) (ok bool) {
	select {
	case o <- element:
		ok = true
	default:
		ok = false
	}
	return
}

// Shift ...
func (o RespQueue) Shift() (element []interface{}, ok bool) {
	select {
	case element = <-o:
		ok = true
	default:
		ok = false
	}
	return
}

// Flush ...
func (o RespQueue) Flush() {
	for {
		_, ok := o.Shift()
		if !ok {
			return
		}
	}
}

// GetLocalRespQueue ...
func GetLocalRespQueue(fullKey string) (o RespQueue) {
	tmp := mygls.Get(fullKey)

	if tmp == nil {
		o = RespQueueNew(queueLength)
		mygls.Set(fullKey, o)
	} else {
		var ok bool
		o, ok = tmp.(RespQueue)
		if !ok {
			panic("methodman.GetLocalRespQueue: the key " + fullKey + " is misused.")
		}
	}
	return
}

// ResetQueue flushed resp queue for a method under current goroutine
func ResetQueue(method interface{}) {
	if !methodtool.IsMethodPointer(method) {
		panic("methodman.MockCleanUp: method is not a method pointer.")
	}
	fullKey := getFullKey(method)
	respQueue := GetLocalRespQueue(fullKey)
	respQueue.Flush()
	return
}

var (
	errorInfoQueueFull = "queue is full. You can review your use case," +
		" or enlarge 'queueLength' by calling SetQueueLength"
)

func getQueueFromMethod(method interface{}) (manager *Manager, respQueue RespQueue) {
	if !methodtool.IsMethodPointer(method) {
		panic("methodman.Expect: method is not a method pointer.")
	}

	methodKey := methodtool.GetMethodUniqueID(method)

	var ok bool
	if manager, ok = managerMap[methodKey]; !ok {
		panic("methodman.Expect: " + errorInfoQueueFull)
	}

	fullKey := getFullKey(reflect.Value(*(manager.Method)).Interface())
	respQueue = GetLocalRespQueue(fullKey)
	return
}
