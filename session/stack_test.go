package session

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stack = "test.stack"

func setupStackRequestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		PushStack(ctx, stack, "A")
		PushStack(ctx, stack, "B")
		PushStack(ctx, stack, "C")

		PushRoute(ctx, "Here")
		PushRoute(ctx, "There")
		ClearRoutes(ctx)
	}
	return http.HandlerFunc(fn)
}

func testStackRequestHandler(t *testing.T) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		assert.Equal(t, "C", PopStack(ctx, stack))
		assert.Equal(t, "B", PopStack(ctx, stack))
		assert.Equal(t, "A", PopStack(ctx, stack))
		assert.Equal(t, "", PopStack(ctx, stack))

		assert.Equal(t, "", PopRoute(ctx))
	}
	return http.HandlerFunc(fn)
}
