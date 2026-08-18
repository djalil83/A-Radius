package genieacs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func GenerateOption43(req Option43Request) (Option43Result, error) {
	if strings.TrimSpace(req.Host) == "" || req.Port < 1 || req.Port > 65535 {
		return Option43Result{}, errors.New("host and valid port are required")
	}
	if req.Protocol != ProtocolHTTP && req.Protocol != ProtocolHTTPS {
		return Option43Result{}, errors.New("unsupported protocol")
	}
	if req.Format != FormatSub01 && req.Format != FormatRaw {
		return Option43Result{}, errors.New("unsupported option 43 format")
	}
	host := req.Host
	if net.ParseIP(host) == nil {
		host = strings.TrimSpace(host)
	}
	acsURL := fmt.Sprintf("%s://%s", strings.ToLower(string(req.Protocol)), net.JoinHostPort(host, strconv.Itoa(req.Port)))
	payload := []byte(acsURL)
	if req.Format == FormatSub01 {
		if len(payload) > 255 {
			return Option43Result{}, errors.New("ACS URL exceeds option 43 sub-option length")
		}
		payload = append([]byte{0x01, byte(len(payload))}, payload...)
	}
	return Option43Result{Hex: "0x" + strings.ToUpper(hex.EncodeToString(payload)), Protocol: req.Protocol, Host: req.Host, Port: req.Port, Username: req.Username}, nil
}

func DecodeOption43(raw string) (Option43Result, error) {
	raw = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(raw, "0x"), "0X"))
	if raw == "" || len(raw)%2 != 0 {
		return Option43Result{}, errors.New("invalid hexadecimal option 43")
	}
	data, err := hex.DecodeString(raw)
	if err != nil {
		return Option43Result{}, errors.New("invalid hexadecimal option 43")
	}
	if len(data) >= 2 && data[0] == 0x01 {
		length := int(data[1])
		if length != len(data)-2 {
			return Option43Result{}, errors.New("sub-option length mismatch")
		}
		data = data[2:]
	}
	u, err := url.Parse(string(data))
	if err != nil || u.Hostname() == "" {
		return Option43Result{}, errors.New("option 43 does not contain a valid ACS URL")
	}
	port := 0
	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return Option43Result{}, errors.New("invalid ACS port")
		}
	}
	protocol := ProtocolHTTP
	if strings.EqualFold(u.Scheme, "https") {
		protocol = ProtocolHTTPS
	} else if !strings.EqualFold(u.Scheme, "http") {
		return Option43Result{}, errors.New("unsupported ACS protocol")
	}
	return Option43Result{Hex: "0x" + strings.ToUpper(hex.EncodeToString(data)), Protocol: protocol, Host: u.Hostname(), Port: port}, nil
}
