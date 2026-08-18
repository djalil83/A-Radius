package customerportal

// CustomerIdentity binds an authenticated APB user to one customer account.
type CustomerIdentity struct {
	UserID     string `json:"user_id"`
	CustomerID string `json:"customer_id"`
}

// CustomerProfile is the safe customer-facing representation.
// Credentials and internal secrets must never be included here.
type CustomerProfile struct {
	ID         string `json:"id"`
	Code       string `json:"customer_code"`
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Address    string `json:"address,omitempty"`
	Village    string `json:"village,omitempty"`
	District   string `json:"district,omitempty"`
	Regency    string `json:"regency,omitempty"`
	Province   string `json:"province,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Status     string `json:"status"`
}
