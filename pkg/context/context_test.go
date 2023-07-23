package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetContextMetadata(t *testing.T) {
	var data = map[ContextMetadata]any{
		USER_ID: "1",
		AUTH:    "token",
	}

	ctx := context.Background()

	ctxVal := SetContext(ctx, data)
	assert.Equal(t, ctxVal.Value(USER_ID), "1")
	assert.Equal(t, ctxVal.Value(AUTH), "token")
}
