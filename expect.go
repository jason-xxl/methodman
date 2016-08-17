package mm

// Expect adds a mocking response to a queue of given method.
// Note that,
// 1) the queue is tied to current goroutine, so call from different goroutine
// won't see this result.
// 2) for the expected response, first in first out
// 3) if no more expected response in queue, the original function would be called
// to serve the call
// 4) you can enqueue mock objects to its constructing function
// 5) it panic if you gave invalid value
//
func Expect(method interface{}, response ...interface{}) {

	_, respQueue := getQueueFromMethod(method)
	ok := respQueue.Push(response)
	if !ok {
		panic("methodman.Expect: queue is full. You can review your use case," +
			" or enlarge 'queueLength' by calling SetQueueLength")
	}
	return
}
