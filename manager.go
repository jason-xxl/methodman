package mm

// Manager maintains states of a method wrapper, which will take over control of
// your original method. The basic idea is, for each calls it will
// 1) check if any expected fake response in the queue for current goroutine
// 2) if yes, consume the response, otherwise call original method like nothing happen
type Manager struct {
	Name   string
	Method *OriginalMethod
	// Queue      RespQueue
}

// ManagerNew ...
func ManagerNew(name string, method interface{}) (o *Manager) {
	originalMethod, ok := OriginalMethodNew(method)
	if !ok {
		panic("methodman.ManagerNew: method param is not a method pointer")
	}
	o = &Manager{
		Name:   name,
		Method: originalMethod,
		//Queue: RespQueueNew(queueLength),
	}
	return o
}
