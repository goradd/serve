package session

import (
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

func TestSetGet(t *testing.T) {
	// setup the ScsSession
	store := memstore.NewWithCleanupInterval(24 * time.Hour)
	sm := scs.New()
	sm.Store = store
	SetSessionManager(NewScsManager(sm))

	// run the session tests
	runRequestTest(t, setRequestHandler(), testRequestHandler(t))
	runRequestTest(t, setupStackRequestHandler(), testStackRequestHandler(t))
}
