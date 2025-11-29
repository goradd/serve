package session

import (
	"context"
)

const key = "goradd.routes"

// PushRoute pushes the given route onto the route stack.
func PushRoute(ctx context.Context, loc string) {
	PushStack(ctx, key, loc)
}

// PopRoute pops the given route off of the route stack and returns it.
func PopRoute(ctx context.Context) (loc string) {
	return PopStack(ctx, key)
}

// ClearRoutes removes all routes from the route stack
func ClearRoutes(ctx context.Context) {
	ClearStack(ctx, key)
}
