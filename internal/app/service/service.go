package service

import "log/slog"

type Services struct {
	Cardcheck
}

func New(log *slog.Logger) *Services {
	return &Services{
		Cardcheck{log: log},
	}
}
