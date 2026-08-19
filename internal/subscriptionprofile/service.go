package subscriptionprofile

import (
	"fmt"
	"regexp"
	"strings"
)

var colorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func NormalizeCreate(req CreateRequest) (CreateRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.ServiceType = strings.ToUpper(strings.TrimSpace(req.ServiceType))
	req.Color = strings.ToLower(strings.TrimSpace(req.Color))
	if req.Color == "" {
		req.Color = "#1677ff"
	}
	if req.SharedUsers == 0 {
		req.SharedUsers = 1
	}
	if req.ActiveDays == 0 {
		req.ActiveDays = 30
	}
	if req.CommissionType == "" {
		req.CommissionType = "RUPIAH"
	}
	if req.BillingCycle == "" {
		req.BillingCycle = "MONTHLY"
	}
	if req.AutoIsolate == nil {
		v := true
		req.AutoIsolate = &v
	}
	if err := validate(req); err != nil {
		return req, err
	}
	return req, nil
}

func NormalizeUpdate(req UpdateRequest) (UpdateRequest, error) {
	base, err := NormalizeCreate(req.CreateRequest)
	req.CreateRequest = base
	if req.Version < 1 {
		return req, fmt.Errorf("version must be >= 1: %w", ErrValidation)
	}
	return req, err
}

func validate(req CreateRequest) error {
	if req.Name == "" || len(req.Name) > 160 {
		return fmt.Errorf("name is required and must be <= 160 characters: %w", ErrValidation)
	}
	switch req.ServiceType {
	case "FTTH", "PPPOE", "HOTSPOT_VOUCHER", "STATIC_IP":
	default:
		return fmt.Errorf("unsupported service_type: %w", ErrValidation)
	}
	if !colorPattern.MatchString(req.Color) {
		return fmt.Errorf("color must be #RRGGBB: %w", ErrValidation)
	}
	if req.SharedUsers < 1 || req.VLANID != nil && (*req.VLANID < 1 || *req.VLANID > 4094) || req.UploadBPS != nil && *req.UploadBPS < 0 || req.DownloadBPS != nil && *req.DownloadBPS < 0 {
		return fmt.Errorf("invalid network values: %w", ErrValidation)
	}
	if req.MonthlyPrice < 0 || req.ActiveDays < 0 || req.CommissionAmount < 0 {
		return fmt.Errorf("invalid billing values: %w", ErrValidation)
	}
	if req.CommissionType != "RUPIAH" && req.CommissionType != "PERCENT" {
		return fmt.Errorf("invalid commission_type: %w", ErrValidation)
	}
	if req.CommissionType == "PERCENT" && req.CommissionAmount > 100 {
		return fmt.Errorf("percentage commission must be <= 100: %w", ErrValidation)
	}
	switch req.BillingCycle {
	case "DAILY", "WEEKLY", "MONTHLY", "CUSTOM":
	default:
		return fmt.Errorf("invalid billing_cycle: %w", ErrValidation)
	}
	for label, value := range map[string]*string{"description": req.Description, "billing_note": req.BillingNote, "rate_limit": req.RateLimit} {
		if value != nil && len(*value) > maxFor(label) {
			return fmt.Errorf("%s is too long: %w", label, ErrValidation)
		}
	}
	return nil
}

func maxFor(label string) int {
	if label == "rate_limit" {
		return 512
	}
	return 2000
}
