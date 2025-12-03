package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/kodra-pay/encryption-service/internal/dto"
)

type EncryptionService struct{}

func NewEncryptionService() *EncryptionService { return &EncryptionService{} }

func (s *EncryptionService) TokenizeCard(_ context.Context, req dto.TokenizeCardRequest) dto.TokenResponse {
	return dto.TokenResponse{
		TokenID: "tok_" + uuid.NewString(),
		Last4:   last4(req.Pan),
		Brand:   "stub",
	}
}

func (s *EncryptionService) GetToken(_ context.Context, id string) dto.TokenResponse {
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
