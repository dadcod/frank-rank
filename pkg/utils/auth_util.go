package utils

import (
	"context"
	"log"

	"github.com/dadcod/frank-rank/internal/middleware"
)

func GetUser(ctx context.Context) *string {
	value := ctx.Value(middleware.AuthUserID)
	log.Println("value")
	log.Println(value)
	if value == nil || value == "" {
		return nil
	}

	userID, ok := value.(string)
	if !ok {
		return nil
	}

	return &userID
}
