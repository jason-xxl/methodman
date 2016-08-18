# methodman
**methodman** is a Go mocking tool based on monkey-patching. 

## Why introduce this tool?

I find most of dependency injection approaches in Go require significant boiletplate codes to get it work. By default it will introdude code structure change (adding dep as extra param) and even logic change, which could potentially make the code much more complex than it should be. And it could be even more complex if the target lib is a 3rdparty libs that doesn't support mocking. So for these reasons I need a tool that can enable mocking with **minimal code change**, and wrote this tool. I like to keep my code as clean / simple / readable as possible.

As extra features, 

- supports parallel unittest. Mocking in one goroutine is invisible in another goroutine.
- supports mocking with a temporary func, which could be useful for simulating timeout, panic, mock with internal state (via closure), or any other kind of side-effects.
- certain helper support to make conversion from integration-test to unit-test easiler. This could be useful in refactoring old code base. See the sections below.

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

Assuming you have a test case that accesses external dependencies and already works fine, assuming no forking gorouting inside the logic, and now, you want to convert it into a unittest. 

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
In test output, you gain the output from original method in the way copy-pastable.
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

## How can I mock an object returned by the dependency pkg? (For refactoring scenarios)

Methodman is modeled around modifiable method, so it won't natively work for this kind of mocking. When converting your implementation to allow Dependency Injection for such case, you would still need to,

1. Abstractise the returned type into an interface, which allow using a mock implementation behind

2. Probably you'll use codegen tool like [Testify Mock](https://github.com/stretchr/testify#mock-package) to generate the mock implementation in independent files

3. By normal practice, you will need to change your main logic to receive dependency as extra param, to allow mocking in unittest.

However, for Step 3, Methodman can make thing easier. You don't really need to change your logic function's signature to add an extra param for receiving dependency. You just need to monkey patch the function where you receive the object, either object constructor, or a singelton getter. Simply let it response with your mock object under the interface in Step 1 in your unittest, and that's all. No need to refactor main logic for mocking.

What about mocking an exported channel without significant refactoring? I'm not sure yet. Please share with me if you got idea.

## How methodman works?

The processing model is simple.

1. When registering the dependency method var, methodman will replace the var with a manager object. It wraps the original method with a queue layer in front (one queue for one method in one goroutine).

2. When you push a fake response to a method, the response will enter the queue of current goroutine.

3. When the method endpoint is called (actually the manager object is called), it will check the queue of current goroutine. If the queue is non-empty, the method will response the fake response by consuming the queue. When queue is empty, the original method is called to provide a real response.

## Complete Demo

Please check out [GitHub Pages](https://github.com/jason-xxl/methodman/blob/master/expect_test.go)
