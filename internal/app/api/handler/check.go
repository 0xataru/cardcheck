package handler

import (
	"cardcheck/internal/domain"
	"log/slog"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Validator interface {
	Validate(card *domain.Card) (*domain.Response, error)
}

type CheckCard struct {
	log          *slog.Logger
	validator    Validator
	reqValidator *validator.Validate
}

// @Summary Validate card
// @Description Validate card - check if card number is valid and expiration date is not in the past
// @Tags check
// @Accept json
// @Produce json
// @Param card body domain.Card true "Card to validate"
// @Success 200 {object} domain.Response
// @Failure 400 {object} domain.ResponseMessage
// @Failure 500 {object} domain.ResponseMessage
// @Router /check [post]
func (cc *CheckCard) Validate(c *fiber.Ctx) error {
	const op = "handler.CheckCard.Validate"
	log := cc.log.With(slog.String("operation", op))

	var card domain.Card

	if err := c.BodyParser(&card); err != nil {
		log.Error("error parsing request body", domain.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseMessage{Message: err.Error()})
	}

	if err := cc.reqValidator.Struct(card); err != nil {
		log.Error("invalid request body", domain.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(domain.ResponseMessage{Message: err.Error()})
	}

	result, err := cc.validator.Validate(&card)
	if err != nil {
		log.Error("error validating card", domain.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ResponseMessage{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
