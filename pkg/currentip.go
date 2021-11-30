package pkg

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type CurrentIP interface {
	Get() (net.IP, error)
}

type ipifyIP struct {
}

type ipifyResponse struct {
	IP string `json:"ip"`
}

func (ipifyIP) Get() (net.IP, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return nil, fmt.Errorf("failed to query ipify: %w", err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoded := ipifyResponse{}
	err = decoder.Decode(&decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response")
	}
	if decoded.IP == "" {
		return nil, fmt.Errorf("could not decode ip")
	}
	ip := net.ParseIP(decoded.IP)
	if ip == nil {
		return nil, fmt.Errorf("ipify return invalid ip")
	}
	return ip, nil
}
