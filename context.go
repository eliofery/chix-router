package chix

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Map шаблон для передачи данных
type Map map[string]any

// Ctx контекст предоставляемый в обработчик
type Ctx struct {
	http.ResponseWriter
	*http.Request
	NextHandler http.Handler

	status int
}

// NewCtx создание контекста
func NewCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		ResponseWriter: w,
		Request:        r,
		NextHandler:    NextHandler(r.Context()),

		status: http.StatusOK,
	}
}

// Status установка статуса ответа
func (ctx *Ctx) Status(status int) *Ctx {
	ctx.status = status
	return ctx
}

// Header установка заголовка
func (ctx *Ctx) Header(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// Decode декодирование тела запроса
func (ctx *Ctx) Decode(data any) error {
	if err := json.NewDecoder(ctx.Request.Body).Decode(data); err != nil {
		ctx.Status(http.StatusBadRequest)

		if errors.Is(err, io.EOF) {
			return errors.New("пустое тело запроса")
		}

		return errors.New("не корректный json")
	}

	return nil
}

// JSON формирование json ответа
func (ctx *Ctx) JSON(data Map) error {
	ctx.Header("Content-Type", "application/json")
	ctx.WriteHeader(ctx.status)

	encoder := json.NewEncoder(ctx.ResponseWriter)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

// Next обработка следующего обработчика
func (ctx *Ctx) Next() error {
	ctx.NextHandler.ServeHTTP(ctx.ResponseWriter, ctx.Request)

	return nil
}
