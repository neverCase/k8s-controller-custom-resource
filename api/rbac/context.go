package rbac

import (
	"context"
	"github.com/Shanghai-Lunara/pkg/jwttoken"
)

const AuthenticationName = "auth"

type Authentication struct {
	TokenClaims *jwttoken.Claims
}

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, u *Authentication) context.Context {
	return context.WithValue(ctx, AuthenticationName, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*Authentication, bool) {
	u, ok := ctx.Value(AuthenticationName).(*Authentication)
	return u, ok
}
