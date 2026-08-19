package customerportal

import "context"

type CustomerDashboard struct {
	Customer CustomerProfile   `json:"customer"`
	Services []CustomerService `json:"services"`
	Summary  DashboardSummary  `json:"summary"`
}

type DashboardSummary struct {
	ServiceCount       int `json:"service_count"`
	ActiveServiceCount int `json:"active_service_count"`
}

func (s *Service) CustomerDashboard(
	ctx context.Context,
	userID string,
) (CustomerDashboard, error) {
	if s == nil || s.repository == nil {
		return CustomerDashboard{}, ErrCustomerNotFound
	}

	customerID, err := s.repository.CustomerIDForUser(ctx, userID)
	if err != nil {
		return CustomerDashboard{}, err
	}

	customer, err := s.repository.GetCustomer(ctx, customerID)
	if err != nil {
		return CustomerDashboard{}, err
	}

	services, err := s.repository.GetServices(ctx, customerID)
	if err != nil {
		return CustomerDashboard{}, err
	}

	active := 0
	for _, service := range services {
		if service.Status == "active" {
			active++
		}
	}

	return CustomerDashboard{
		Customer: customer,
		Services: services,
		Summary: DashboardSummary{
			ServiceCount:       len(services),
			ActiveServiceCount: active,
		},
	}, nil
}
