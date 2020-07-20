package payments

import (
	"context"
	"fmt"
)

// Storage defines user service's storage interface.
type Storage interface {
	SourceTypes(context.Context) ([]SourceType, error)
	ProceedPayment(context.Context, Payment) error
	Balance(context.Context) (Account, error)
}

// Service implements user functionality.
type Service struct {
	storage     Storage
	sourceTypes map[string]int
}

// NewService returns a new instance of Service.
func NewService(storage Storage) (*Service, error) {
	sourceTypes, err := storage.SourceTypes(context.Background())
	if err != nil {
		return nil, err
	}

	s := Service{
		storage:     storage,
		sourceTypes: make(map[string]int),
	}

	for _, v := range sourceTypes {
		s.sourceTypes[v.Value] = v.ID
	}

	return &s, nil
}

// SourceTypeID returns all source types
func (s *Service) SourceTypeID(_ context.Context, sourceType string) (int, error) {
	id, ok := s.sourceTypes[sourceType]
	if !ok {
		return 0, fmt.Errorf("wrong header Source-Type")
	}

	return id, nil
}

// ProceedPayment processes payments
func (s *Service) ProceedPayment(ctx context.Context, payment Payment) error {
	if err := s.storage.ProceedPayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to proceed payment: %w", err)
	}
	return nil
}

// Balance returns account balance
func (s *Service) Balance(ctx context.Context) (Account, error) {
	return s.storage.Balance(ctx)
}
