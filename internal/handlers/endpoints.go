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
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid token ID")
	}
	return c.JSON(h.svc.GetToken(c.Context(), id))
}

// EncryptData encrypts arbitrary data
func (h *EncryptionHandler) EncryptData(c *fiber.Ctx) error {
	var req struct {
		Data string `json:"data"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Data == "" {
		return fiber.NewError(fiber.StatusBadRequest, "data field is required")
	}

	encrypted, err := h.svc.EncryptData(c.Context(), req.Data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "encryption failed")
	}

	return c.JSON(fiber.Map{
		"encrypted_data": encrypted,
	})
}

// DecryptData decrypts encrypted data
func (h *EncryptionHandler) DecryptData(c *fiber.Ctx) error {
	var req struct {
		EncryptedData string `json:"encrypted_data"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.EncryptedData == "" {
		return fiber.NewError(fiber.StatusBadRequest, "encrypted_data field is required")
	}

	decrypted, err := h.svc.DecryptData(c.Context(), req.EncryptedData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "decryption failed")
	}

	return c.JSON(fiber.Map{
		"data": decrypted,
	})
}
