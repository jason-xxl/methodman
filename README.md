# methodman

**methodman** is a mocking tool for Go based on monkey-patching. 

In many dynamic language it's easy to monkey patching object or methods. However in Go, exported pkg method is not modifiable. So to get monkey patching work, one (and only) prerequition is define the target method in a way like this (in the case you can control the code), 
```
  var TheMethodToBeMocked = func(...){...}
```  
or define a reference var in caller side like,
```
  var DependencyRef = targetpkg.TargetMethod
```  
Then you can mock on DependencyRef.

### Why introduce this tool?

I find most of dependency injection approaches in Go require significant boiletplate codes to get it work. This usually will introdude code structure change and even logic change, which could sometimes make the code significantly complex than it should be. For some cases it's even more complex if the target lib is a 3rdparty libs that's not built with allowing mocking in mind. So for these reasons I need a tool that can enable mocking with **minimal code footprint**, and wrote this tool. I like to keep my code as clean / readable as possible.

As extra features, methodman is equipped with [GoroutineLocalStorage](https://github.com/tylerb/gls), which is used to enables parallel unittest, that mocking in one goroutine won't mocking in another, and it also support mocking with a temporary func, which is useful to simulate timeout or various side-effects.

### Examples

Please check out [GitHub Pages](https://github.com/jason-xxl/methodman/blob/master/expect_test.go)
