package entities

import (
	"github.com/gofiber/fiber/v2"
	"github.com/korvised/go-ecommerce/pkg/logger"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, traceId, msg string) IResponse
	Res() error
}

type Response struct {
	StatusCode int
	Data       any
	ErrorRes   *ErrorResponse
	context    *fiber.Ctx
	IsError    bool
}

type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"messages"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	logger.InitLogger(r.context, &r.Data).Print().Save()

	return r
}

func (r *Response) Error(code int, traceId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: traceId,
		Msg:     msg,
	}
	r.IsError = true
	logger.InitLogger(r.context, &r.ErrorRes).Print().Save()

	return r
}

func (r *Response) Res() error {
	return r.context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes
		}
		return &r.Data
	}())
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
