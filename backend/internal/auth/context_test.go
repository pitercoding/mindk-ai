package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIDFromContext_Present(t *testing.T) {
	ctx := WithUserID(context.Background(), "user_abc")

	userID, ok := UserIDFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, "user_abc", userID)
}

func TestUserIDFromContext_Missing(t *testing.T) {
	userID, ok := UserIDFromContext(context.Background())

	assert.False(t, ok)
	assert.Empty(t, userID)
}
