package cms

import "context"

type PageDefaults interface {
	GetDefaults(ctx context.Context) (map[string]any, error)
	GetRouteDefaults(ctx context.Context, routeName string) (map[string]any, error)
}
