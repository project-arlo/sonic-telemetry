package common_utils

import (
	"context"
	"fmt"
	"sync/atomic"
)


// AuthInfo holds data about the authenticated user
type AuthInfo struct {
	// Username
	User string
	AuthEnabled bool
	// Roles
	Roles []string
}

// RequestContext holds metadata about REST request.
type RequestContext struct {

	// Unique reqiest id
	ID string

	// Auth contains the authorized user information
	Auth AuthInfo
}

type contextkey int

const requestContextKey contextkey = 0

// Request Id generator
var requestCounter uint64

// GetContext function returns the RequestContext object for a
// HTTP request. RequestContext is maintained as a context value of
// the request. Creates a new RequestContext object is not already
// available; in which case this function also creates a copy of
// the HTTP request object with new context.
func GetContext(ctx context.Context) (*RequestContext, context.Context) {
	cv := ctx.Value(requestContextKey)
	if cv != nil {
		return cv.(*RequestContext), ctx
	}

	rc := new(RequestContext)
	rc.ID = fmt.Sprintf("TELEMETRY-%v", atomic.AddUint64(&requestCounter, 1))

	ctx = context.WithValue(ctx, requestContextKey, rc)
	return rc, ctx
}
