package subscriptionprofile

import "errors"

var (
	ErrNotFound       = errors.New("subscription profile not found")
	ErrConflict       = errors.New("subscription profile version conflict")
	ErrValidation     = errors.New("subscription profile validation failed")
	ErrInvalidRequest = errors.New("invalid request")
)

type Profile struct {
	ID               string  `json:"id"`
	TenantID         string  `json:"tenant_id"`
	Name             string  `json:"name"`
	ServiceType      string  `json:"service_type"`
	Category         *string `json:"category,omitempty"`
	Media            *string `json:"media,omitempty"`
	Color            string  `json:"color"`
	Description      *string `json:"description,omitempty"`
	Status           string  `json:"status"`
	MikrotikGroup    *string `json:"mikrotik_group,omitempty"`
	RadiusGroup      *string `json:"radius_group,omitempty"`
	RateLimit        *string `json:"rate_limit,omitempty"`
	UploadBPS        *int64  `json:"upload_bps,omitempty"`
	DownloadBPS      *int64  `json:"download_bps,omitempty"`
	SharedUsers      int     `json:"shared_users"`
	VLANID           *int    `json:"vlan_id,omitempty"`
	OLTProfile       *string `json:"olt_profile,omitempty"`
	IPPool           *string `json:"ip_pool,omitempty"`
	MonthlyPrice     int64   `json:"monthly_price"`
	ActiveDays       int     `json:"active_days"`
	CommissionAmount int64   `json:"commission_amount"`
	CommissionType   string  `json:"commission_type"`
	BillingCycle     string  `json:"billing_cycle"`
	AutoIsolate      bool    `json:"auto_isolate"`
	BillingNote      *string `json:"billing_note,omitempty"`
	Version          int64   `json:"version"`
	CreatedBy        *string `json:"created_by,omitempty"`
	UpdatedBy        *string `json:"updated_by,omitempty"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type CreateRequest struct {
	Name             string  `json:"name"`
	ServiceType      string  `json:"service_type"`
	Category         *string `json:"category"`
	Media            *string `json:"media"`
	Color            string  `json:"color"`
	Description      *string `json:"description"`
	MikrotikGroup    *string `json:"mikrotik_group"`
	RadiusGroup      *string `json:"radius_group"`
	RateLimit        *string `json:"rate_limit"`
	UploadBPS        *int64  `json:"upload_bps"`
	DownloadBPS      *int64  `json:"download_bps"`
	SharedUsers      int     `json:"shared_users"`
	VLANID           *int    `json:"vlan_id"`
	OLTProfile       *string `json:"olt_profile"`
	IPPool           *string `json:"ip_pool"`
	MonthlyPrice     int64   `json:"monthly_price"`
	ActiveDays       int     `json:"active_days"`
	CommissionAmount int64   `json:"commission_amount"`
	CommissionType   string  `json:"commission_type"`
	BillingCycle     string  `json:"billing_cycle"`
	AutoIsolate      *bool   `json:"auto_isolate"`
	BillingNote      *string `json:"billing_note"`
}

type UpdateRequest struct {
	CreateRequest
	Version int64 `json:"version"`
}

type ListResult struct {
	Items  []Profile `json:"items"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type Revision struct {
	ID        int64   `json:"id"`
	ProfileID string  `json:"profile_id"`
	Version   int64   `json:"version"`
	Operation string  `json:"operation"`
	Snapshot  []byte  `json:"snapshot"`
	ChangedBy *string `json:"changed_by,omitempty"`
	ChangedAt string  `json:"changed_at"`
}
