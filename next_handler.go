package chix

import (
	"context"
	"net/http"
)

// Тип key ключ контекста
type key string

// Значение ключа
const nextKey key = "next"

// WithNextHandler добавление следующего обработчика в цепочке middleware в контекст
func WithNextHandler(ctx context.Context, next http.Handler) context.Context {
	return context.WithValue(ctx, nextKey, next)
}

// NextHandler получение следующего обработчика в цепочке middleware из контекста
func NextHandler(ctx context.Context) http.Handler {
	val := ctx.Value(nextKey)

	next, ok := val.(http.Handler)
	if !ok {
		return nil
	}

	return next
}
