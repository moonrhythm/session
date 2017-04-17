package session

import (
	"context"

	"github.com/acoshift/middleware"
)

type contextType int

const (
	contextTypeSession contextType = iota
)

// Middleware is the session parser middleware
func Middleware() middleware.Middleware {
	return nil
}

// Session type
type Session struct {
	UserID interface{}
}

// Get gets session from context
func Get(ctx context.Context) *Session {
	sess, _ := ctx.Value(contextTypeSession).(*Session)
	return sess
}
