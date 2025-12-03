package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/encryption-service/internal/handlers"
	"github.com/kodra-pay/encryption-service/internal/services"
)

func Register(app *fiber.App, service string) {
	health := handlers.NewHealthHandler(service)
	health.Register(app)

	svc := services.NewEncryptionService()
	h := handlers.NewEncryptionHandler(svc)
	api := app.Group("/")
	api.Post("tokenize/card", h.TokenizeCard)
	api.Get("tokens/:id", h.GetToken)
}
