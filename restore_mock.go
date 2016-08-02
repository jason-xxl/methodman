package methodman

import (
	mygls "github.com/tylerb/gls"
)

// RestoreMock will detach queue internal states to allow garbage collection.
// It's required for every goroutine that calls Expect. Usually it's set as a defer
// in the top level function of the gouroutine.
func RestoreMock() {
	mygls.Cleanup()
}
