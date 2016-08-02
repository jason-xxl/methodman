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

##### Firstly, attach a method agent to the method to be mocked
```
func TestMain(m *testing.M) {
	flag.Parse()
	EnableMock(&dep_pkg.MethodA, "MethodA")
	os.Exit(m.Run())
}
```
##### Then, mock it in your test
```
func TestNormalUse(t *testing.T) {

	defer RestoreMock()
	Expect(&dep_pkg.MethodA, "some fake response for my test as 1st returned var", "some more, as 2nd retuened var")

	// then you can receive above 2 value in your code path.
	// if all fake responses are consumed, the agent will fall back to original method.
	ret1, ret 2 := dep_pkg.MethodA(1, "2")
}
```

### Complete Demo

Please check out [GitHub Pages](https://github.com/jason-xxl/methodman/blob/master/expect_test.go)
