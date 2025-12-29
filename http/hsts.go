package http

import (
	"fmt"
	"net/http"
)

// WriteHstsHeader writes an HSTS header to w, which will cause the browser to enforce HTTPS by redirecting
// http requests to port 443 over https.
//
// maxAge is in seconds. A value of zero will disable hsts.
// includeSubDomains will cause attempts at contacting subdomains over http to also redirect to https.
// preload signals intent to be included in Google's special universal preload domain list.
func WriteHstsHeader(w http.ResponseWriter, maxAge int64, includeSubDomains bool, preload bool) {
	if preload {
		includeSubDomains = true // required
		// google requires a one-year minimum maxAge, but we are not going to enforce that here, since it might change.
	}
	out := fmt.Sprintf("max-age=%d", maxAge)
	if includeSubDomains {
		out += "; includeSubDomains"
	}
	if preload {
		out += "; preload"
	}
	w.Header().Set("Strict-Transport-Security", out)
}
