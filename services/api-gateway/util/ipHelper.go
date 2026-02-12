package util

import (
	"encoding/json"
	"net/http"
)

func CountryByIP(ip string) (string, error) {
	resp, err := http.Get("https://ipapi.co/" + ip + "/json/")
	if err != nil {
		return "0.0.0.0", err
	}
	defer resp.Body.Close()

	var r IpResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "0.0.0.0", err
	}

	return r.CountryCode, nil
}
