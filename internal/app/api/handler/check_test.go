package handler

import (
	"bytes"
	"cardcheck/internal/domain"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func Test_CheckCard_Validate(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	reqValidator := validator.New()

	tests := []struct {
		name           string
		card           *domain.Card
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid request",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "12",
				ExpirationYear:  "2024",
			},
			expectedStatus: fiber.StatusOK,
			expectedError:  "",
		},
		{
			name: "no card number",
			card: &domain.Card{
				CardNumber:      "",
				ExpirationMonth: "12",
				ExpirationYear:  "2024",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "CardNumber",
		},
		{
			name: "no expiration month",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "",
				ExpirationYear:  "2024",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "ExpirationMonth",
		},
		{
			name: "no expiration year",
			card: &domain.Card{
				CardNumber:      "5167803252097675",
				ExpirationMonth: "12",
				ExpirationYear:  "",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "ExpirationYear",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			cc := &CheckCard{
				log: log,
				validator: mockValidator{
					card: tt.card,
				},
				reqValidator: reqValidator,
			}

			app.Post("/check", cc.Validate)

			body, _ := json.Marshal(tt.card)
			req := httptest.NewRequest("POST", "/check", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			responseBody, _ := io.ReadAll(resp.Body)
			responseStr := string(responseBody)

			if tt.expectedError != "" {
				if !strings.Contains(strings.ToLower(responseStr), strings.ToLower(tt.expectedError)) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.expectedError, responseStr)
				}
			} else {
				expectedSuccess := `{"valid":true,"error":null}`
				if responseStr != expectedSuccess {
					t.Errorf("expected success response '%s', got '%s'", expectedSuccess, responseStr)
				}
			}
		})
	}
}

type mockValidator struct {
	card *domain.Card
}

func (v mockValidator) Validate(card *domain.Card) (*domain.Response, error) {
	if card.CardNumber == v.card.CardNumber {
		return &domain.Response{
			Valid: true,
			Error: nil,
		}, nil
	}

	return nil, errors.New("Invalid card number")
}
