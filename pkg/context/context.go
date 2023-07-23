package context

import (
	"context"
)

type ContextMetadata int

const (
	AUTH ContextMetadata = iota + 1
	USER_ID
	PROFILE_ID
	PHONE
	EMAIL
	STATUS
)

func SetContext(ctx context.Context, list map[ContextMetadata]any) context.Context {
	for index, val := range list {
		ctx = context.WithValue(ctx, index, val)
	}

	return ctx
}
