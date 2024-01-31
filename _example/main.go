package main

import (
	"fmt"
	"github.com/eliofery/go-chix"
	"net/http"
	"time"
)

type Request struct {
	Name string
	Age  int
}

type Response struct {
	Date time.Time
}

func main() {
	route := chix.NewRouter()

	route.Get("/profile", func(ctx *chix.Ctx) error {
		var req Request
		if err := ctx.Decode(&req); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(chix.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		var res Response
		res.Date = time.Now()
		return ctx.JSON(chix.Map{
			"success": true,
			"message": "Время ответа",
			"data":    res,
		})
	})

	route.With(Cors(), Example()).Route("/group", func(r *chix.Router) {
		r.Get("/route1", func(ctx *chix.Ctx) error { return nil })
		r.Get("/route2", func(ctx *chix.Ctx) error { return nil })
	})

	route.Group(func(r *chix.Router) {
		r.Use(Example())

		r.Get("/route1", func(ctx *chix.Ctx) error { return nil })
		r.Get("/route2", func(ctx *chix.Ctx) error { return nil })
	})

	if err := route.Listen("127.0.0.1:3000"); err != nil {
		fmt.Println(err)
	}
}
