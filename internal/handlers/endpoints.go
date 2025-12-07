package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/encryption-service/internal/dto"
	"github.com/kodra-pay/encryption-service/internal/services"
)

type EncryptionHandler struct {
	svc *services.EncryptionService
}

func NewEncryptionHandler(svc *services.EncryptionService) *EncryptionHandler {
	return &EncryptionHandler{svc: svc}
}

func (h *EncryptionHandler) TokenizeCard(c *fiber.Ctx) error {
	var req dto.TokenizeCardRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.TokenizeCard(c.Context(), req))
}

func (h *EncryptionHandler) GetToken(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id") // Use c.ParamsInt
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid token ID")
	}
	return c.JSON(h.svc.GetToken(c.Context(), id))
}
