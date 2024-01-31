# Chix роутер

Chix роутер для Golang позволяет создавать обработчики с контекстом, используя роутер Chi.

Обработчик маршрутов как в фреймворках chi, fiber и т.п.

## Пример стандартного обработчика Chi:

```go
route.Get("/", func Handler(w http.ResponseWrite, r *http.Request) {})

route.Use(func(next http.Handler) http.Handler {})
```

## Пример обработчика Chix:

```go
route.Get("/", func(ctx *chix.Ctx) error {})

route.Use(func(ctx *chix.Ctx) error {})
```

## Реализованные методы:

- Get
- Post
- Put
- Patch
- Delete
- NotFound
- MethodNotAllowed
- Use
- Group
- With
- Route
- Mount
- ServeHTTP
- Listen
- 