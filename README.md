# methodman

**methodman** is a Go mocking tool based on monkey-patching. 

## How to enable monkey patching in Golang?

In many dynamic languages it's easy to monkey patching object or methods. However, in Go, exported pkg method is not modifiable. So by default there's no formal way to monkey-patch. To get monkey patching works, pkg method have to be defined as an exported method variable to allow modify (in case you can control the code),
```
  var TheMethodToBeMocked = func(...){...}
```  
Or, use a reference var in caller side like this (in case you can't control the code),
```
  var TheMethodToBeMocked = targetpkg.TargetMethod
```  
Then you can monkey patch TheMethodToBeMocked for mocking.

Given the target method var is modifiable, methodman would replace TheMethodToBeMocked with a wrapper, who provides mockability that just overlays the original method. 

## Why introduce this tool?

I find most of dependency injection approaches in Go require significant boiletplate codes to get it work. By default it will introdude code structure change (adding dep as extra param) and even logic change, which could potentially make the code much more complex than it should be. And it could be even more complex if the target lib is a 3rdparty libs that doesn't support mocking. So for these reasons I need a tool that can enable mocking with **minimal code change**, and wrote this tool. I like to keep my code as clean / simple / readable as possible.

As extra features, 

- supports parallel unittest (assuming no side effect) by equipting [GoroutineLocalStorage](https://github.com/tylerb/gls). Mocking in one goroutine won't affect mocking in another goroutine.
- supports mocking with a temporary func, which could be useful for simulating timeout, panic, mock with internal state (via closure), or any other kind of side-effects.
- certain helper support to make conversion from integration-test to unit-test easiler. This could be useful in refactoring old code base. See the sections below.

## How to use?

##### 1. Install.
```
go get -u github.com/jason-xxl/methodman
```
##### 2. Register the method to be mocked with a valid name.
```
func TestMain(m *testing.M) {
	flag.Parse()
	mm.EnableMock(&dep_pkg.MethodA, "dep_pkg.MethodA")
	os.Exit(m.Run())
}
```
##### 3. Mock it in your test.
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

## How to convert integration-test into unittest? (For refactoring scenarios)

Assuming you have a test case that accesses external dependencies and already works fine, assuming no forking gorouting insight the logic, now, you want to convert it into unittest. 

So the idea here is you run your integration-test once, capture those real outputs of depending methods, and use those output as sample to mock. Here's step by step.

##### 1. Enable the `CapturingLogger` and register the methods that you want to mock. 
```
func TestMain(m *testing.M) {
	flag.Parse()

	mm.SetLogger(mm.CapturingLogger)

	mm.EnableMock(&dep_pkg.MethodA, "dep_pkg1.MethodA")
	mm.EnableMock(&dep_pkg.MethodB, "dep_pkg2.MethodA")

	os.Exit(m.Run())
}
```
##### 2. Run your test in verbose mode, you can get the real response in the form of usable code that you can insert into your code. This save human effort to form the mock response and reduce human error.
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
##### 3. By copying them into your test, you gain the unittest version based on previous integration.
```
func TestToBeConverted(t *testing.T) {

    defer mm.Init(t).CleanUp()

    mm.Expect(&dep_pkg1.MethodA, "real response 1", "real response 2")
    mm.Expect(&dep_pkg2.MethodB, "real response 3")

    ... // Your original code here. No change.
}
```
##### 4. After that, remove this line.
```
mm.SetLogger(mm.CapturingLogger)
```

Now you got a perfect unittest.

## Complete Demo

Please check out [GitHub Pages](https://github.com/jason-xxl/methodman/blob/master/expect_test.go)
