package customerportal

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CustomerProfile(
	ctx context.Context,
	userID string,
) (CustomerProfile, error) {
	if s == nil || s.repository == nil {
		return CustomerProfile{}, ErrCustomerNotFound
	}

	customerID, err := s.repository.CustomerIDForUser(ctx, userID)
	if err != nil {
		return CustomerProfile{}, err
	}

	return s.repository.GetCustomer(ctx, customerID)
}

func (s *Service) CustomerServices(
	ctx context.Context,
	userID string,
) ([]CustomerService, error) {
	if s == nil || s.repository == nil {
		return nil, ErrCustomerNotFound
	}

	customerID, err := s.repository.CustomerIDForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.repository.GetServices(ctx, customerID)
}
