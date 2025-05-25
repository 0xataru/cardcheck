package handler

import (
	"log/slog"

	"github.com/go-playground/validator"
)

type ServiceInterface interface {
	Validator
}

type Handler struct {
	CheckCard
}

func New(s ServiceInterface, v *validator.Validate, l *slog.Logger) *Handler {
	return &Handler{CheckCard: CheckCard{log: l, validator: s, reqValidator: v}}
}
