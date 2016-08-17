package mm

import (
	"strings"

	"github.com/kr/pretty"

	"fmt"
	"sync"
)

// Logger should be an external method for logging the captured content
// by default the logger is empty and capturing is disabled
// note: Assuming it's set once at testing beginning, it's not Mutex protected.
// Don't change it during test.
type Logger func(methodName string, isMockResponse bool, output []interface{})

// SetLogger Logger would be used during the test, default is disabled,
// you should set it up (only once) before test runs if you need it
func SetLogger(logger Logger) {
	currentLoggerOnce.Do(func() {
		currentLogger = logger
	})
}

var (
	sourceRealOutput = "real"

	currentLogger     = Logger(nil)
	currentLoggerOnce sync.Once

	loggerMutex = sync.Mutex{}

	// DefaultLogger is for minimal and human readable logging
	DefaultLogger = func(methodName string, isMockResponse bool, output []interface{}) {
		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		var tag string
		if isMockResponse {
			tag = "[mock]"
		} else {
			tag = "[real]"
		}

		tpl := tag + " %s"

		t := ", %#v"
		t = strings.Repeat(t, len(output))

		tpl += t
		tpl += ""

		output = append([]interface{}{methodName}, output...)

		code := pretty.Sprintf(tpl, output...)

		fmt.Println(code)
	}

	// CapturingLogger organise the real output in code form that you can copy paste.
	// It's a handy helper for converting integration test into unittest.
	// 1. set it as logger
	// 2. call "go test -v -run yourtest"
	CapturingLogger = func(methodName string, isMockResponse bool, output []interface{}) {
		if !isMockResponse {
			return
		}

		loggerMutex.Lock()
		defer loggerMutex.Unlock()

		tpl := "mm.Expect(&%s"

		t := ", %#v"
		t = strings.Repeat(t, len(output))

		tpl += t
		tpl += ")"

		output = append([]interface{}{methodName}, output...)

		code := pretty.Sprintf(tpl, output...)

		fmt.Println(code)
	}
)
