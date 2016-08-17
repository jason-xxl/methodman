package mm

import (
	"testing"

	"strconv"
	"sync"

	"flag"
	"os"

	"time"
)

var (
	MethodAInfo = "I'm original MethodA"

	// So, MethodA here represents an exported method from some library you want to mock
	// By convention, it's in a form of method variable
	//
	MethodA = func(p1 int, p2 string) (r1 string) {
		return MethodAInfo
	}
)

func TestMain(m *testing.M) {
	flag.Parse()

	// For this demo, we know we need to mock certain external
	// method to give us fake response, so we'll need to upgrade
	// that method "MethodA". To avoid race condition, it's
	// required to be put at the beginning of all tests.
	//
	// We firstly check the original method key, which is unique
	// based on its method body pointer
	//
	methodKeyA0 := GetMethodUniqueID(MethodA)

	// By calling EnableMock, a method manager will wrap the original method and
	// switch the method var of exported method to the method manager itself.
	//
	// If we have multiple methods to upgrade for all the unittest here, we
	// need to upgrade all of them in TestMain.
	//
	// From this time on, all method calls are trapped into method manager.
	// If there's a fake resp in the queue of current goroutine, by calling
	// MethodA, the test will receive the fake resp, otherwise, the original
	// method body would be executed.
	//
	EnableMock(&MethodA, "MethodA")

	// After that we can check the current method key.
	//
	// See? It's already the key of the method manager, not the
	// original method.
	//
	methodKeyA1 := GetMethodUniqueID(MethodA)
	if methodKeyA0 == methodKeyA1 {
		panic("EnableMock should have replaced original method")
	}

	// (Optionally) You can set a custom logger to trace what's responsed.
	//
	// For cases you want to convert an existing integration to unittest, you can
	// use a predefined CapturingLogger which would output response from original
	// methods as usable code, in calling order. (go test -v .)
	//
	// SetLogger(CapturingLogger)

	// ok, initiation done. We can start real unittests.
	//
	os.Exit(m.Run())
}

func TestNormalUse(t *testing.T) {

	// Setup cleanup as defer. It should be added to each goroutine
	// where Expect is called. Usually as defer
	//
	defer Init(t).CleanUp()

	// So now we can inject some fake resp into MethodA endpoint.
	//
	// Note that, if the format of fake resp does not match the method,
	// the test would finally panic by the mismatching, panic by reflect pkg.
	// For this demo we skipped this panic test first.
	//
	info1 := "Now I'm method manager of MethodA, not original MethodA"
	Expect(&MethodA, info1)

	// So calling same method var would receive the mocked value above
	// by consuming it from the queue behind.
	//
	// ** Importantly, the queue is only visible to current goroutine.
	// ** Different goroutine will own different queue (per method).
	// ** So it supports parallel unittest naturally. If there's no fake
	// ** response remain in queue, calling method manager of MethodA
	// ** should be exactly same with calling original method body of
	// ** MethodA.
	//
	ret1 := MethodA(1, "2")
	if ret1 != info1 {
		t.Fatal("Manager should response queued resp. " + ret1)
	}

	// Additionally, it supports temporary implementation of the original method
	// This is good for some special test case, e.g. simulating a timeout.
	info3 := "I'm temporary implementation."
	ExpectFunc(&MethodA, func(p1 int, p2 string) (r1 string) {

		// simulate timeout
		time.Sleep(1 * time.Nanosecond)
		return info3
	})

	// Calling it will see the 100 ms delay before response.
	//
	ret3 := MethodA(1, "2")
	if ret3 != info3 {
		t.Fatal("Manager should sleep and response. " + info3)
	}

	// Calling it again. Because now the queue is empty, the manager
	// is falling back by calling the original method.
	//
	ret4 := MethodA(1, "2")
	if ret4 != MethodAInfo {
		t.Fatal("Manager should call original method when queued empty. " + ret4)
	}

	// Now let's go with multiple (10k) goroutines and each 3 fake resp
	//
	var wg sync.WaitGroup
	for idx := 0; idx < 10000; idx++ {

		wg.Add(1)

		go func(i int) {

			// Setup cleanup as defer. It should be added to each goroutine
			// where Expect is called. Usually as defer
			//
			defer Init(t).CleanUp()

			info1 := "I'm method manager of MethodA " + strconv.Itoa(i)

			// test 3 times
			Expect(&MethodA, info1)
			Expect(&MethodA, info1)
			Expect(&MethodA, info1)

			ret1 := MethodA(1, "2")
			if ret1 != info1 {
				t.Fatal("Manager should response queued resp. " + ret1)
			}

			ret2 := MethodA(1, "2")
			if ret2 != info1 {
				t.Fatal("Manager should response queued resp. " + ret2)
			}

			ret3 := MethodA(1, "2")
			if ret3 != info1 {
				t.Fatal("Manager should response queued resp. " + ret3)
			}

			// falling back to original method body
			ret4 := MethodA(1, "2")
			if ret4 != MethodAInfo {
				t.Fatal("Manager should call original method when queued empty. " + ret4)
			}

			wg.Done()
		}(idx)
	}

	// So by this test you'll see that unittest on one goroutine is fully isolated
	// from those on another. Now you can feel free to run a bunch of unrelated unittests
	// in parallel mode with arbitrary mocking that won't disturb each other.
	//
	wg.Done()
}
