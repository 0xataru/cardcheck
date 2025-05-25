package domain

import (
	"log/slog"
)

type ErrorCode string

const (
	InvalidCardNumber     ErrorCode = "001"
	InvalidExpirationDate ErrorCode = "002"
)

type Error struct {
	Code    ErrorCode `json:"code" example:"001"`
	Message string    `json:"message" example:"error message"`
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
