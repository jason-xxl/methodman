# methodman

**methodman** is a Go mocking tool based on monkey-patching. 

In many dynamic language it's easy to monkey patching object or methods. However in Go, exported pkg method is not modifiable. So to get monkey patching work, pkg method need to be defined as a method variable to allow monkey patching (in the case you can control the code),
```
  var TheMethodToBeMocked = func(...){...}
```  
or use a reference var in caller side like this (for cases you can't change code of the lib),
```
  var TheMethodToBeMocked = targetpkg.TargetMethod
```  
then you can monkey patch TheMethodToBeMocked for mocking.

Given the method var is modifiable, methodman would replace TheMethodToBeMocked with a wrapper, who provides mockability that overlays the original method. 

### Why introduce this tool?

I find most of dependency injection approaches in Go require significant boiletplate codes to get it work. This usually introdudes code structure change and even logic change, which could sometimes make the code significantly more complex than it should be. For some cases it's even more complex if the target lib is a 3rdparty libs that's not built with allowing mocking in mind. So for these reasons I need a tool that can enable mocking with **minimal code footprint**, and wrote this tool. I like to keep my code as clean / readable as possible.

As extra features, methodman is equipped with [GoroutineLocalStorage](https://github.com/tylerb/gls), which is used to enables parallel mocking / unittest, that mocking in one goroutine won't affect mocking in another goroutine. It also supports mocking with a temporary func, which could be useful for simulating timeout, panic, mock with internal state (via closure), or any other kind of side-effects.

### How to use

##### Install
```
go get -u github.com/jason-xxl/methodman
```
##### Firstly, attach a method agent to the method to be mocked
```
func TestMain(m *testing.M) {
	flag.Parse()
	mm.EnableMock(&dep_pkg.MethodA, "MethodA")
	os.Exit(m.Run())
}
```
##### Then, mock it in your test
```
func TestNormalUse(t *testing.T) {

	defer mm.Init(t).CleanUp()
	
	mm.Expect(&dep_pkg.MethodA, "some fake response as 1st returned var", "some more, as 2nd retuened var")

	// Then you can receive above 2 value in your code path.
	ret1, ret 2 := dep_pkg.MethodA(1, "2")
	
	// 1. If all fake responses are consumed, the agent will fall back to original method.
	// 2. It doesn't matter next call of dep_pkg.MethodA is at which level, above fake value would be 
	//    received if it's in same goroutine.
}
```

### Converting Integration Test to Unittest

Assuming you have a test case that accesses external dependencies and already works find, you want to convert into unittest.

1. Enable the `CapturingLogger` and register the methods that you want to mock.
```
func TestMain(m *testing.M) {
	flag.Parse()

	mm.SetLogger(mm.CapturingLogger)

	mm.EnableMock(&dep_pkg.MethodA, "dep_pkg1.MethodA")
	mm.EnableMock(&dep_pkg.MethodB, "dep_pkg2.MethodA")

	os.Exit(m.Run())
}
```
2. Run your test in verbose mode, you can get the real response in the form of usable code that you can insert into your code. This save human effort to form the mock response and reduce human error.
```
go test -v -run TestToBeConverted
```
You would gain output from original method in the way copy-pastable.
```
...
mm.Expect(&dep_pkg1.MethodA, "real response 1", "real response 2")
mm.Expect(&dep_pkg2.MethodB, "real response 3")
...
...
```
3. By copying them into your test, you gain the unittest version based on previous integration.
```
func TestToBeConverted(t *testing.T) {

    defer mm.Init(t).CleanUp()

    mm.Expect(&dep_pkg1.MethodA, "real response 1", "real response 2")
    mm.Expect(&dep_pkg2.MethodB, "real response 3")

    ... // Your original code here. No change.
}
```
4. After that, remove this line.
```
mm.SetLogger(mm.CapturingLogger)
```

Now you got a perfect unittest.


### Complete Demo

Please check out [GitHub Pages](https://github.com/jason-xxl/methodman/blob/master/expect_test.go)
