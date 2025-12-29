package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/goradd/serve/log"
	strings2 "github.com/goradd/strings"
)

// ParseValueAndParams returns the value and param map for Content-Type and Content-Disposition header values.
func ParseValueAndParams(in string) (value string, params map[string]string) {
	parts := strings.Split(in, ";")
	if len(parts) > 0 {
		value = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			for _, p := range parts[1:] {
				p = strings.TrimSpace(p)
				offset := strings.IndexRune(p, '=')
				if offset >= 0 {
					if params == nil {
						params = make(map[string]string)
					}
					params[p[:offset]] = p[offset+1:]
				}
			}
		}
	}
	return
}

// ParseAuthorizationHeader will parse an authorization header into its scheme and params.
func ParseAuthorizationHeader(auth string) (scheme, params string) {
	var found bool
	before, after, found := strings.Cut(auth, " ")
	scheme = before
	if found {
		params = strings.TrimSpace(after)
	}
	return
}

// ValidateHeader confirms that the given header's values only contains ASCII characters.
func ValidateHeader(header http.Header) bool {
	for k, a := range header {
		if !strings2.IsASCII(k) {
			log.Info(nil, logModule, "A header key did not contain only ASCII values",
				slog.String("key", k))
			return false
		}
		for _, h := range a {
			if !strings2.IsASCII(h) {
				log.Info(nil, logModule, "A header value did not contain only ASCII values",
					slog.String("key", k),
					slog.String("value", h))
				return false
			}
		}
	}
	return true
}

// WithHeaderValidator insert middleware into the handler stack that performs OWASP style validation on a request.
func WithHeaderValidator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !ValidateHeader(r.Header) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
