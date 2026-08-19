package subscriptionprofile

import (
	"strings"
	"testing"
)

func TestNormalizeCreateDefaults(t *testing.T) {
	req, err := NormalizeCreate(CreateRequest{Name: " Rumahan ", ServiceType: "ftth"})
	if err != nil {
		t.Fatal(err)
	}
	if req.Name != "Rumahan" || req.ServiceType != "FTTH" || req.Color != "#1677ff" || req.SharedUsers != 1 || req.ActiveDays != 30 || req.AutoIsolate == nil || !*req.AutoIsolate {
		t.Fatalf("unexpected defaults: %+v", req)
	}
}

func TestNormalizeCreateRejectsInvalidValues(t *testing.T) {
	cases := []CreateRequest{
		{Name: "", ServiceType: "FTTH"},
		{Name: "x", ServiceType: "UNKNOWN"},
		{Name: "x", ServiceType: "FTTH", Color: "blue"},
		{Name: "x", ServiceType: "FTTH", SharedUsers: 0 - 1},
		{Name: "x", ServiceType: "FTTH", CommissionType: "PERCENT", CommissionAmount: 101},
	}
	for i, tc := range cases {
		if _, err := NormalizeCreate(tc); err == nil {
			t.Errorf("case %d expected error", i)
		}
	}
}

func TestNormalizeUpdateRequiresVersion(t *testing.T) {
	_, err := NormalizeUpdate(UpdateRequest{CreateRequest: CreateRequest{Name: "x", ServiceType: "FTTH"}})
	if err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("expected version error, got %v", err)
	}
}
