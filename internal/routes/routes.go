package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/encryption-service/internal/handlers"
	"github.com/kodra-pay/encryption-service/internal/services"
)

func Register(app *fiber.App, serviceName string, encryptionKey string) {
	// Health check
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	// Initialize encryption service
	encSvc, err := services.NewEncryptionService(encryptionKey)
	if err != nil {
		panic("failed to initialize encryption service: " + err.Error())
	}

	// Encryption handlers
	encHandler := handlers.NewEncryptionHandler(encSvc)

	// API routes
	api := app.Group("/api/v1")

	// Token management
	api.Post("/keys/tokenize", encHandler.TokenizeCard)
	api.Get("/keys/:id", encHandler.GetToken)

	// Data encryption/decryption
	api.Post("/encrypt", encHandler.EncryptData)
	api.Post("/decrypt", encHandler.DecryptData)
}
