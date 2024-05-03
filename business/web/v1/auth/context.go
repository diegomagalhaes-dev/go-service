package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const claimKey ctxKey = 1

const userKey ctxKey = 2

func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

func GetClaims(ctx context.Context) Claims {
	v, ok := ctx.Value(claimKey).(Claims)
	if !ok {
		return Claims{}
	}
	return v
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userKey, userID)
}

func GetUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return v
}
