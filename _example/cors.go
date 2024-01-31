package main

import (
	"github.com/eliofery/go-chix"
	"github.com/go-chi/cors"
)

const defaultCorsMaxAge = 3600 // 1 час

// Cors настройки межсайтового взаимодействия
// Пример: https://github.com/go-chi/cors?tab=readme-ov-file#usage
func Cors() chix.Handler {
	return func(ctx *chix.Ctx) error {
		corsHandler := cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link", "Content-Length", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
			MaxAge:           defaultCorsMaxAge,
		})

		corsHandler(ctx.NextHandler).ServeHTTP(ctx.ResponseWriter, ctx.Request)

		return nil
	}
}
