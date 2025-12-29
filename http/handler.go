package http

import (
	"context"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/goradd/serve/config"
)

// The following two maps collect handler registration during Go's init process. These are
// then registered to the application muxers when the application starts up. This makes it
// possible for parts of the app to turn themselves on just by being imported.

// PatternMuxer is the muxer that immediately routes handlers based on the path without
// going through the application handlers. It responds directly to the ResponseWriter.
var PatternMuxer Muxer = http.NewServeMux()

// AppMuxer is the application muxer that lets you do http handling
// from behind the application facilities of session management, output buffering, etc.
var AppMuxer Muxer = http.NewServeMux()

// RegisterStaticHandler registers a handler for the given pattern.
//
// The given handler is served immediately by the application without going through the application
// handler stack. If you need session management, HSTS protection, authentication, etc., use
// RegisterAppHandler.
//
// If a ProxyPath is set, it will automatically be inserted in front of the path in the pattern.
// If the pattern has a path that ends in "/"
func RegisterStaticHandler(pattern string, handler http.Handler) {
	pattern = joinProxyPath(pattern)
	PatternMuxer.Handle(pattern, handler)
}

// RegisterAppHandler registers a handler for the given pattern.
//
// Use this when registering a handler to a specific path. Use RegisterAppPrefixHandler if registering
// a handler for a whole subdirectory of a path.
//
// The given handler is served near the end of the application handler stack, so
// you will have access to session management and any other middleware handlers
// in the application stack.
//
// You may call this from an init() function.
func RegisterAppHandler(pattern string, handler http.Handler) {
	pattern = joinProxyPath(pattern)
	AppMuxer.Handle(pattern, handler)
}

// A DrawFunc sends output to the Writer. goradd uses this signature in its template functions.
type DrawFunc func(ctx context.Context, w io.Writer) (err error)

// RegisterDrawFunc registers an output function for the given pattern.
//
// This could be used to register template output with a path, for example. See the renderResource
// template macro and the configure.tpl.got file in the welcome application for an example.
//
// The file name extension will be used first to determine the Content-Type. If that fails, then
// the content will be inspected to determine the Content-Type.
//
// Registered handlers are served by the AppMuxer.
func RegisterDrawFunc(pattern string, f DrawFunc) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		name := path.Base(r.URL.Path)

		// set the content type by extension first
		ctype := mime.TypeByExtension(filepath.Ext(name))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}

		// the write process below will fill in the content-type if not set
		err := f(ctx, w)
		if err != nil {
			panic(err)
		}
	}
	h := http.HandlerFunc(fn)
	RegisterAppHandler(pattern, h)
}

func joinProxyPath(pattern string) string {
	if config.ProxyPath == "" {
		return pattern
	}

	// assume the path is well-formed
	offset := strings.IndexRune(pattern, '/')
	newPattern := pattern[:offset] + config.ProxyPath + pattern[offset:]
	return newPattern
}

// WithAppMuxer serves up the AppMuxer.
func WithAppMuxer(next http.Handler) http.Handler {
	return WithMuxer(AppMuxer, next)
}

// WithPatternMuxer serves up the PatternMuxer.
func WithPatternMuxer(next http.Handler) http.Handler {
	return WithMuxer(PatternMuxer, next)
}
