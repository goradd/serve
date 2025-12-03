package serve

import (
	"fmt"
	"net/http"

	"github.com/goradd/base"
	"github.com/goradd/serve/config"
	http2 "github.com/goradd/serve/http"
)

// ServerBaseI defines the virtual functions that are callable on the Server.
type ServerBaseI interface {
	Init()
	DefaultHandler() http.Handler
	HSTSHandler(next http.Handler) http.Handler
	ServeAppMux(next http.Handler) http.Handler
	ServePatternMux(next http.Handler) http.Handler

	/*
		ServeHTTP(w http.ResponseWriter, r *http.Request)
		PutContext(*http.Request) *http.Request
		SetupErrorHandling()
		SetupPagestateCaching()
		SetupSessionManager()
		SetupMessenger()
		SetupDatabaseWatcher()
		SetupPaths()
		SessionHandler(next http.Handler) http.Handler
		AccessLogHandler(next http.Handler) http.Handler
		PutDbContextHandler(next http.Handler) http.Handler
		ServeRequestHandler() http.Handler

	*/
}

type ServerBase struct {
	base.Base
	//httpErrorReporter http2.ErrorReporter
}

func (a *ServerBase) Init(self ServerBaseI) {
	a.Base.Init(self)
}

func (a *ServerBase) this() ServerBaseI {
	return a.Self().(ServerBaseI)
}

func (a *ServerBase) MakeHandler() http.Handler {
	// the handler chain gets built in the reverse order of getting called

	// These handlers are called in reverse order
	h := a.this().DefaultHandler() // Should go at the end of the chain to catch whatever is missed above
	h = a.this().ServeAppMux(h)    // Serves other dynamic files, and possibly the api
	//	h = a.ServePageHandler(h)           // Serves the Goradd dynamic pages
	//	h = a.PutAppContextHandler(h)
	//	h = a.this().SessionHandler(h)
	//	h = a.BufferedOutputHandler(h) // Must be in front of the session handler
	//	h = a.StatsHandler(h)
	h = a.this().ServePatternMux(h) // Serves most static files and websocket requests.
	// Must be after the error handler so panics are intercepted by the error reporter
	// and must be in front of the buffered output handler because of websocket server
	//	h = a.this().PutDbContextHandler(h) // This is here so that the PatternMux handlers can use the ORM
	//	h = a.validateHttpHandler(h)
	//	h = a.httpErrorReporter.Use(h) // Default http error handler to intercept panics.
	h = a.this().HSTSHandler(h)
	//	h = a.this().AccessLogHandler(h)

	return h
}

// DefaultHandler is the last handler on the default call chain.
// It returns a simple not found error.
// Note that the html root handler registered in embedder.go also handles situations where an http
// path is not found.
// You can override this handler by duplicating it in your app object.
func (a *ServerBase) DefaultHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
	return http.HandlerFunc(fn)
}

// HSTSHandler sets the browser to HSTS mode using the given timeout. HSTS will force a browser to accept only
// HTTPS connections for everything coming from your domain, if the initial page was served over HTTPS. Many browsers
// already do this. What this additionally does is prevent the user from overriding this. Also, if your
// certificate is bad or expired, it will NOT allow the user the option of using your website anyway.
// This should be safe to send in development mode if your local server is not using HTTPS, since the header
// is ignored if a page is served over HTTP.
//
// Once the HSTS policy has been sent to the browser, it will remember it for the amount of time
// specified, even if the header is not sent again. However, you can override it by sending another header, and
// clear it by setting the timeout to 0. Set the timeout to -1 to turn it off. You can also completely override this by
// implementing this function in your app.go file.
func (a *ServerBase) HSTSHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if config.HSTSTimeout >= 0 {
			w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains", config.HSTSTimeout))
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ServeAppMux serves up the http2.AppMuxer, which handles REST calls,
// and dynamically created files.
//
// To use your own AppMuxer, override this function in app.go and create your own.
func (a *ServerBase) ServeAppMux(next http.Handler) http.Handler {
	return http2.UseMuxer(http2.AppMuxer, next)
}

// ServePatternMux serves up the http2.PatternMuxer.
//
// The pattern muxer routes patterns early in the handler stack. It is primarily for handlers that
// do not need the session manager or buffered output, things like static files.
//
// The default version injects a standard http muxer. Override to use your own muxer.
func (a *ServerBase) ServePatternMux(next http.Handler) http.Handler {
	return http2.UseMuxer(http2.PatternMuxer, next)
}
