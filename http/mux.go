package http

import (
	"net/http"
)

// Muxer represents the typical functions available in a mux and allows you
// to replace the default muxer here with a 3rd party mux, like the Gorilla mux.
//
// However, beware. The default Go muxer will do redirects. If this goradd application
// is behind a reverse proxy that is rewriting the url, the Go muxer will not correctly
// do rewrites because it will not include the reverse proxy path in the rewrite
// rule, and things will break.
//
// If you create your own mux and you want to do redirects, use MakeLocalPath to
// create the redirect url. See also maps.SafeMap for a map you can use if you
// are modifying paths while using the mux.
type Muxer interface {
	// Handle associates a handler with the given pattern in the url path
	Handle(pattern string, handler http.Handler)

	// Handler returns the handler associate with the request, if one exists. It
	// also returns the actual path registered to the handler
	Handler(r *http.Request) (h http.Handler, pattern string)

	// ServeHTTP sends a request to the MUX, to be forwarded on to the registered handler,
	// or responded with an unknown resource error.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// UseMuxer serves a muxer such that if a handler cannot be found, or the found handler does not respond,
// control is past to the next handler.
//
// Note that the default Go Muxer is NOT recommended, as it improperly handles
// redirects if this is behind a reverse proxy.
func UseMuxer(mux Muxer, next http.Handler) http.Handler {
	if next == nil {
		panic("next may not be nil. Pass a http.NotFoundHandler if this is the end of the handler chain")
	}
	if mux == nil {
		panic("mux may not be nil")
	}
	fn := func(w http.ResponseWriter, r *http.Request) {

		var h http.Handler
		var p string

		h, p = mux.Handler(r)
		if p == "" {
			// not found, so go to next handler
			next.ServeHTTP(w, r) // skip
		} else {
			// match, so serve normally
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
