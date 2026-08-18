package genieacs

import "time"

type DeviceStatus string

const (
	Online         DeviceStatus = "ONLINE"
	OfflineUnder24 DeviceStatus = "OFFLINE_UNDER_24H"
	OfflineOver24  DeviceStatus = "OFFLINE_OVER_24H"
	Unknown        DeviceStatus = "UNKNOWN"
)

// CalculateDeviceStatus derives status from backend timestamps; frontend input is never trusted.
func CalculateDeviceStatus(now time.Time, lastInform, lastConnected *time.Time) DeviceStatus {
	last := lastInform
	if last == nil || (lastConnected != nil && lastConnected.After(*last)) {
		last = lastConnected
	}
	if last == nil {
		return Unknown
	}
	age := now.Sub(*last)
	if age < 5*time.Minute {
		return Online
	}
	if age < 24*time.Hour {
		return OfflineUnder24
	}
	return OfflineOver24
}

type Device struct {
	ID         string       `json:"id"`
	CustomerID string       `json:"customer_id"`
	BranchID   string       `json:"branch_id"`
	Username   string       `json:"username,omitempty"`
	PPPoE      string       `json:"pppoe_ip,omitempty"`
	TR069      string       `json:"tr069_ip,omitempty"`
	Serial     string       `json:"serial_number"`
	PON        string       `json:"pon_port,omitempty"`
	Vendor     string       `json:"manufacturer,omitempty"`
	Model      string       `json:"model,omitempty"`
	Firmware   string       `json:"firmware,omitempty"`
	SSID       string       `json:"ssid,omitempty"`
	WANStatus  string       `json:"wan_status"`
	Status     DeviceStatus `json:"status"`
	Server     string       `json:"server,omitempty"`
	LastSeen   *time.Time   `json:"last_seen,omitempty"`
	LastInform *time.Time   `json:"last_inform_at,omitempty"`
}

type Server struct {
	ID               string  `json:"id"`
	BranchID         string  `json:"branch_id"`
	Name             string  `json:"name"`
	Host             string  `json:"host"`
	Port             int     `json:"port"`
	Username         string  `json:"username"`
	Password         *string `json:"password"`
	CredentialRef    string  `json:"credential_ref,omitempty"`
	Status           string  `json:"status"`
	ConnectionStatus string  `json:"connection_status"`
}

type Command string

const (
	CommandSummon Command = "SUMMON"
	CommandReboot Command = "REBOOT"
	CommandDelete Command = "DELETE"
	CommandReset  Command = "RESET"
	CommandSync   Command = "SYNC"
	CommandDHCP43 Command = "DHCP_OPTION_43"
)

var ApprovalRequiredCommands = map[Command]bool{
	CommandSummon: true, CommandReboot: true, CommandDelete: true,
	CommandReset: true, CommandSync: true, CommandDHCP43: true,
}

var HighRiskCommands = map[Command]bool{CommandDelete: true, CommandReset: true}

type CommandProposal struct {
	ID                string   `json:"id"`
	Command           Command  `json:"command"`
	DeviceIDs         []string `json:"device_ids"`
	TargetCount       int      `json:"target_count"`
	Reason            string   `json:"reason"`
	Status            string   `json:"status"`
	ApprovalRequired  bool     `json:"approval_required"`
	ProductionChanged bool     `json:"production_changed"`
	RequestedBy       string   `json:"requested_by"`
	ApprovedBy        *string  `json:"approved_by,omitempty"`
}

type Option43Protocol string
type Option43Format string

const (
	ProtocolHTTP  Option43Protocol = "HTTP"
	ProtocolHTTPS Option43Protocol = "HTTPS"
	FormatSub01   Option43Format   = "SUBOPTION_01_LENGTH"
	FormatRaw     Option43Format   = "RAW"
)

type Option43Request struct {
	Host       string           `json:"host"`
	Port       int              `json:"port"`
	Username   string           `json:"username,omitempty"`
	Password   string           `json:"password,omitempty"`
	Protocol   Option43Protocol `json:"protocol"`
	Format     Option43Format   `json:"format"`
	OptionName string           `json:"option_name"`
}

type Option43Result struct {
	Hex      string           `json:"hex"`
	Protocol Option43Protocol `json:"protocol"`
	Host     string           `json:"host"`
	Port     int              `json:"port"`
	Username string           `json:"username,omitempty"`
}
