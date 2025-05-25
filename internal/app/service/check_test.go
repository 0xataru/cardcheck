package service

import (
	"cardcheck/internal/domain"
	"log/slog"
	"os"
	"testing"
)

func Test_Validate(t *testing.T) {
	cc := &Cardcheck{
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}

	tests := []struct {
		name      string
		card      *domain.Card
		wantValid bool
		wantErr   *domain.Error
	}{
		{
			name: "valid card",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "12",
				ExpirationYear:  "2029",
			},
			wantValid: true,
			wantErr:   nil,
		},
		{
			name: "invalid card number - too short",
			card: &domain.Card{
				CardNumber:      "1234",
				ExpirationMonth: "12",
				ExpirationYear:  "2024",
			},
			wantValid: false,
			wantErr: &domain.Error{
				Code:    domain.InvalidCardNumber,
				Message: "invalid card number",
			},
		},
		{
			name: "invalid card number - fails Luhn check",
			card: &domain.Card{
				CardNumber:      "5457626723237072",
				ExpirationMonth: "12",
				ExpirationYear:  "2024",
			},
			wantValid: false,
			wantErr: &domain.Error{
				Code:    domain.InvalidCardNumber,
				Message: "invalid card number",
			},
		},
		{
			name: "invalid expiration month - too short",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "1",
				ExpirationYear:  "2024",
			},
			wantValid: false,
			wantErr: &domain.Error{
				Code:    domain.InvalidExpirationDate,
				Message: "invalid expiration date",
			},
		},
		{
			name: "invalid expiration month - too large",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "13",
				ExpirationYear:  "2024",
			},
			wantValid: false,
			wantErr: &domain.Error{
				Code:    domain.InvalidExpirationDate,
				Message: "invalid expiration date",
			},
		},
		{
			name: "invalid expiration year - in the past",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "12",
				ExpirationYear:  "2020",
			},
			wantValid: false,
			wantErr: &domain.Error{
				Code:    domain.InvalidExpirationDate,
				Message: "invalid expiration date",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cc.Validate(tt.card)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if tt.wantErr == nil {
				if result.Error != nil {
					t.Errorf("error = %v, want nil", result.Error)
				}
			} else {
				if result.Error == nil {
					t.Errorf("error = nil, want %v", tt.wantErr)
				} else if result.Error.Code != tt.wantErr.Code {
					t.Errorf("error code = %s, want %s", result.Error.Code, tt.wantErr.Code)
				}
			}
		})
	}
}
