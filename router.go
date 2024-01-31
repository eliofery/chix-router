package chix

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Handler обработчик
type Handler func(ctx *Ctx) error

// Router обертка над chi роутером
type Router struct {
	*chi.Mux
}

// NewRouter создание роутера
func NewRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

// handleCtx запускает обработчик роутера
func (rt *Router) handler(handler Handler, w http.ResponseWriter, r *http.Request) {
	ctx := NewCtx(w, r)

	if err := handler(ctx); err != nil {
		err = ctx.JSON(Map{
			"success": false,
			"message": err.Error(),
		})
		if err != nil {
			http.Error(ctx.ResponseWriter, "Не предвиденная ошибка", http.StatusInternalServerError)
		}
	}
}

// Get запрос на получение данных
func (rt *Router) Get(path string, handler Handler) {
	rt.Mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Post запрос на сохранение данных
func (rt *Router) Post(path string, handler Handler) {
	rt.Mux.Post(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Put запрос на обновление всех данных
func (rt *Router) Put(path string, handler Handler) {
	rt.Mux.Put(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Patch запрос на обновление конкретных данных
func (rt *Router) Patch(path string, handler Handler) {
	rt.Mux.Patch(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Delete запрос на удаление данных
func (rt *Router) Delete(path string, handler Handler) {
	rt.Mux.Delete(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// NotFound обрабатывает 404 ошибку
func (rt *Router) NotFound(handler Handler) {
	rt.Mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// MethodNotAllowed обрабатывает 405 ошибку
func (rt *Router) MethodNotAllowed(handler Handler) {
	rt.Mux.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rt.handler(handler, w, r)
	})
}

// Use добавляет промежуточное программное обеспечение
func (rt *Router) Use(middlewares ...Handler) {
	for _, middleware := range middlewares {
		currentMiddleware := middleware

		rt.Mux.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				ctx = WithNextHandler(ctx, next)

				rt.handler(currentMiddleware, w, r.WithContext(ctx))
			})
		})
	}
}

// With добавляет встроенное промежуточное программное обеспечение для обработчика конечной точки
func (rt *Router) With(middlewares ...Handler) *Router {
	var handlers []func(http.Handler) http.Handler

	for _, middleware := range middlewares {
		currentMiddleware := middleware

		handler := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				ctx = WithNextHandler(ctx, next)

				rt.handler(currentMiddleware, w, r.WithContext(ctx))
			})
		}

		handlers = append(handlers, handler)
	}

	return &Router{
		Mux: rt.Mux.With(handlers...).(*chi.Mux),
	}
}

// Route создает вложенность роутеров
func (rt *Router) Route(pattern string, fn func(r *Router)) *Router {
	subRouter := &Router{
		Mux: chi.NewRouter(),
	}

	fn(subRouter)
	rt.Mount(pattern, subRouter)

	return subRouter
}

// Mount добавляет вложенность роутеров
func (rt *Router) Mount(pattern string, router *Router) {
	rt.Mux.Mount(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.Mux.ServeHTTP(w, r)
	}))
}

// Group группирует роутеры
func (rt *Router) Group(fn func(r *Router)) *Router {
	im := rt.With()

	if fn != nil {
		fn(im)
	}

	return im
}

// ServeHTTP возвращает весь пул роутеров
func (rt *Router) ServeHTTP() http.HandlerFunc {
	// Здесь ни чего особенного просто возвращаем стандартный http.HandlerFunc
	return rt.Mux.ServeHTTP
}

// Listen запускает сервер
// Реализация: https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
func (rt *Router) Listen(addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: rt.ServeHTTP(),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Не удалось запустить сервер: %s", err.Error())
				ch <- ctx.Err()
			}
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeoutCtx, done := context.WithTimeout(context.Background(), time.Second*10)
		defer done()

		go func() {
			<-timeoutCtx.Done()
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				fmt.Printf("Время корректного завершения работы истекло. Принудительный выход: %s", timeoutCtx.Err().Error())
			}
		}()

		if err := server.Shutdown(timeoutCtx); err != nil {
			fmt.Printf("Не удалось остановить сервер: %s", err.Error())
		}
	}

	return nil
}
