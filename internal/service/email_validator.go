package service

import (
	"context"

	"github.com/GoSeoTaxi/email-validator/internal/domain"
	"github.com/GoSeoTaxi/email-validator/internal/pb"
)

type EmailValidatorService struct {
	pb.UnimplementedEmailValidatorServer
	Validator domain.EmailValidator
}

func NewEmailValidatorService(validator domain.EmailValidator) *EmailValidatorService {
	return &EmailValidatorService{
		Validator: validator,
	}
}

func (s *EmailValidatorService) ValidateEmail(ctx context.Context, req *pb.ValidateEmailRequest) (*pb.ValidateEmailResponse, error) {
	_ = ctx
	isValid, message := s.Validator.Validate(req.Email)
	return &pb.ValidateEmailResponse{
		IsValid: isValid,
		Message: message,
	}, nil
}
