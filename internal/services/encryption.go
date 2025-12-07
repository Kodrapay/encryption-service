package services

import (
	"context"

	"github.com/kodra-pay/encryption-service/internal/dto"
)

type EncryptionService struct{}

func NewEncryptionService() *EncryptionService { return &EncryptionService{} }

func (s *EncryptionService) TokenizeCard(_ context.Context, req dto.TokenizeCardRequest) dto.TokenResponse {
	// In a real scenario, this would generate a unique int ID.
	// For this mock implementation, we return a placeholder.
	// req.Cvv (int) and req.Reference (int) are now available in the request.
	return dto.TokenResponse{
		TokenID: 1, // Placeholder for an auto-generated int ID
		Last4:   last4(req.Pan),
		Brand:   "stub",
	}
}

func (s *EncryptionService) GetToken(_ context.Context, id int) dto.TokenResponse {
	return dto.TokenResponse{
		TokenID: id,
		Last4:   "0000",
		Brand:   "stub",
	}
}

func last4(pan string) string {
	if len(pan) < 4 {
		return pan
	}
	return pan[len(pan)-4:]
}
