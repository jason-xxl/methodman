# methodman

*methodman* is a mocking tool for Go based on monkey-patching. 

In many dynamic language it's easy to monkey patching object or methods. However in Go, exported pkg method is not modifiable. So to get monkey patching work, one (and only) prerequition is define the target method in a way like this (in the case you can control the code), 

  var TheMethodToBeMocked = func(...){...}
  
or define a reference var in caller side like,

  var DependencyRef = targetpkg.TargetMethod
  
Then you can mock on DependencyRef.

# Why introduce this tool?

I find most of dependency injection in Go require a lot boiletplate code to get it work, this usually will introdude code structure change and logic change, which could sometimes make the code significantly complex than it should be. For some cases it's even more complex if the target lib is a 3rdparty lib that's not built with Go Dependency Injection in mind. So for these cases I need a tool that can enable mocking with minimal code footprint.

As extra features, the tool is equipped with GLS from tylerb, which enables parallel unittest, that mocking in one goroutine won't mocking in another, and it also support mocking with a temporary func, which is useful to simulate timeout or various side-effects.

