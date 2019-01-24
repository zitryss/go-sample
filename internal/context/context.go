package context

import (
	"context"
	"time"
)

type ctxKey int

const (
	createdOn ctxKey = iota
	createdBy
)

func WithCreatedOn(ctx context.Context, val time.Time) context.Context {
	return context.WithValue(ctx, createdOn, val)
}

func CreatedOn(ctx context.Context) time.Time {
	return ctx.Value(createdOn).(time.Time)
}

func WithCreatedBy(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, createdBy, val)
}

func CreatedBy(ctx context.Context) string {
	return ctx.Value(createdBy).(string)
}
