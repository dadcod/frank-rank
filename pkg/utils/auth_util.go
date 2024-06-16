package utils

import (
	"context"

	"github.com/dadcod/frank-rank/internal/middleware"
)

func GetUser(ctx context.Context) *string {
	value := ctx.Value(middleware.AuthUserID)
	if value == nil || value == "" {
		return nil
	}

	userID, ok := value.(string)
	if !ok {
		return nil
	}

	return &userID
}
