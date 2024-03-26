package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/rarime-auth-svc/resources"
	"github.com/rarimo/rarime-link-svc/internal/data"
	points "github.com/rarimo/rarime-points-svc/pkg/connector"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	storageCtxKey
	userClaimCtxKey
	pointsCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxStorage(storage data.Storage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, storageCtxKey, storage)
	}
}

func Storage(r *http.Request) data.Storage {
	return r.Context().Value(storageCtxKey).(data.Storage)
}

func CtxUserClaim(claim []resources.Claim) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userClaimCtxKey, claim)
	}
}

func UserClaim(r *http.Request) []resources.Claim {
	return r.Context().Value(userClaimCtxKey).([]resources.Claim)
}

func CtxPoints(pointsCon *points.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, pointsCtxKey, pointsCon)
	}
}

func Points(r *http.Request) *points.Client {
	return r.Context().Value(pointsCtxKey).(*points.Client)
}
