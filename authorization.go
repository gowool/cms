package cms

import (
	"context"

	"github.com/gowool/cms/model"
)

type AuthScheme int8

const (
	UnknownScheme AuthScheme = iota
	BasicScheme
	JWTScheme
)

type Decision int8

const (
	DecisionDeny = iota + 1
	DecisionAllow
)

func (d Decision) String() string {
	switch d {
	case DecisionDeny:
		return "deny"
	case DecisionAllow:
		return "allow"
	default:
		return "unknown"
	}
}

type Access int8

const (
	AccessUnknown Access = iota
	AccessPublic
	AccessPrivate
	AccessRead
	AccessWrite
	AccessAdmin
)

type Claims struct {
	Subject  *model.Admin
	Scheme   AuthScheme
	TwoFA    bool
	Metadata map[string]any
}

type Decider interface {
	Decide(ctx context.Context, claims *Claims) (Decision, error)
}

func NewDecider(access Access, twoFA bool) Decider {
	return NewTargetAccess(access, twoFA)
}

type TargetAccess struct {
	Access Access
	TwoFA  bool
}

func NewTargetAccess(access Access, twoFA bool) *TargetAccess {
	return &TargetAccess{Access: access, TwoFA: twoFA}
}

func (a *TargetAccess) Decide(_ context.Context, claims *Claims) (Decision, error) {
	if claims == nil {
		return DecisionDeny, nil
	}
	if a.TwoFA && !claims.TwoFA {
		return DecisionDeny, nil
	}
	if a.Access == AccessUnknown {
		return DecisionDeny, nil
	}

	role := model.RoleGuest
	if claims.Subject != nil {
		role = claims.Subject.Role
	}

	if role >= getRequiredRole(a.Access) {
		return DecisionAllow, nil
	}
	return DecisionDeny, nil
}

type CallTarget struct {
	OperationID string
	Access      map[AuthScheme]Decider
	Metadata    map[string]any
}

func NewCallTarget(access Access) *CallTarget {
	return &CallTarget{
		Access: map[AuthScheme]Decider{
			BasicScheme: NewDecider(access, false),
			JWTScheme:   NewDecider(access, true),
		},
	}
}

type Authorizer interface {
	Authorize(ctx context.Context, claims *Claims, target *CallTarget) (Decision, error)
}

type DefaultAuthorizer struct{}

func NewDefaultAuthorizer() *DefaultAuthorizer {
	return &DefaultAuthorizer{}
}

func (*DefaultAuthorizer) Authorize(ctx context.Context, claims *Claims, target *CallTarget) (Decision, error) {
	if claims == nil || target == nil || target.Access == nil {
		return DecisionDeny, nil
	}

	decider, ok := target.Access[claims.Scheme]
	if !ok {
		return DecisionDeny, nil
	}
	return decider.Decide(ctx, claims)
}

func getRequiredRole(access Access) model.Role {
	switch access {
	case AccessPublic:
		return model.RoleGuest
	case AccessPrivate, AccessRead:
		return model.RoleReader
	case AccessWrite:
		return model.RoleWriter
	default:
		return model.RoleAdmin
	}
}
