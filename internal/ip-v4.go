package internal

import (
	"io"
	"net/http"
)

type IPv4 struct {
	Address string
}

func (u *IPv4) MetricName() string {
	return "ipv4"
}

func GetIPv4() (*IPv4, error) {
	resp, err := http.Get("https://api.ipify.org")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &IPv4{
		Address: string(body),
	}, nil
}
