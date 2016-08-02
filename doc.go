// Package methodman tries to minimise code change to support Dependency Injection
// to support mocking in unittest. The assumption it makes is very minimal, just
// make sure the exported methods from your lib are in a form of exported variables.
// So it provides a manager for your method behind the scene, if you enqueued a
// fake response for you test in current goroutine, calling the (managed) method
// will firstly response your fake response instead to calling the original.
// Note that, for putting minimal code footprint as one of top priority, this lib will
// use panic for failed assertion check (like parameter check) instead of passing err
// back to user that simply make final code bloated.
package methodman
