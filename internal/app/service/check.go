package service

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"cardcheck/internal/domain"
)

type Cardcheck struct {
	log *slog.Logger
}

func (cc *Cardcheck) Validate(card *domain.Card) (*domain.Response, error) {
	const op = "service.Cardcheck.Validate"
	log := cc.log.With(slog.String("operation", op))
	log.Info("attempting to validate card")

	response := domain.Response{
		Valid: true,
		Error: nil,
	}

	expMonth, err := strconv.Atoi(card.ExpirationMonth)
	if err != nil {
		log.Warn("invalid expiration month")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	expYear, err := strconv.Atoi(card.ExpirationYear)
	if err != nil {
		log.Warn("invalid expiration year")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !isValidCardNumber(card.CardNumber) {
		log.Warn("invalid card number")
		response.Valid = false
		response.Error = &domain.Error{
			Code:    domain.InvalidCardNumber,
			Message: "invalid card number",
		}
		return &response, nil
	}

	if !isValidExpirationDate(expMonth, expYear) {
		log.Warn("invalid expiration date")

		response.Valid = false
		response.Error = &domain.Error{
			Code:    domain.InvalidExpirationDate,
			Message: "invalid expiration date",
		}
		return &response, nil
	}

	log.Info("card has been validated")

	return &response, nil
}

func isValidCardNumber(cardNumber string) bool {
	// Check length first as it's the fastest check
	if len(cardNumber) < 12 || len(cardNumber) > 19 {
		return false
	}

	// Check if all characters are digits
	if _, err := strconv.Atoi(cardNumber); err != nil {
		return false
	}

	// Luhn algorithm:
	//
	//validates card numbers by checking if the sum of digits
	// (with special doubling) is divisible by 10
	//
	// https://en.wikipedia.org/wiki/Luhn_algorithm

	sum := 0
	isSecondDigit := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecondDigit = !isSecondDigit
	}

	return sum%10 == 0
}

func isValidExpirationDate(expirationMonth, expirationYear int) bool {
	currentYear, currentMonth, _ := time.Now().Date()

	// Check if month is valid (1-12)
	if expirationMonth < 1 || expirationMonth > 12 {
		return false
	}

	// Check if year is in the past
	if expirationYear < int(currentYear) {
		return false
	}

	// Check if current month is past expiration month in current year
	if expirationYear == int(currentYear) && expirationMonth < int(currentMonth) {
		return false
	}

	return true
}
