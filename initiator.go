package mm

import (
	"testing"

	"github.com/tylerb/gls"
)

var (
	// defaultInitiator ...
	defaultInitiator = &Initiator{
		CleanUp: func() {
			gls.Cleanup()
		},
	}

	// queueKeyPrefixT is used to identify currentT in goroutine local storage
	queueKeyPrefixT = "_methodman_queue_132435_*testing.T"
)

type Initiator struct {
	// CleanUp will do clean up gls state for current goroutine
	CleanUp func()
}

// Init store the current *testing.T for use
// By default you can use it at beginning of your test like this way
//     defer mm.Init(t).CleanUp()
func Init(t *testing.T) *Initiator {
	if t == nil {
		panic("methodman.Init: t is nil")
	}
	gls.Set(queueKeyPrefix+"*testing.T", t)
	return defaultInitiator
}

// GetCurrentT gets current T object from gls
func GetCurrentT() (t *testing.T) {
	o := gls.Get(queueKeyPrefixT)
	if o != nil {
		t, _ = o.(*testing.T)
	}
	return
}
