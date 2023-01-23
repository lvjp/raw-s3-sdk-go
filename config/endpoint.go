package config

import (
	"fmt"
	neturl "net/url"
	"strconv"
)

type Endpoint struct {
	Host string
	Port int

	WithSSL         bool
	WithVirtualHost bool
}

func NewEndpointFromURL(url string) (e Endpoint, err error) {
	u, err := neturl.Parse(url)
	if err != nil {
		err = fmt.Errorf("cannot parse testserver URL: %w", err)
		return
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		err = fmt.Errorf("cannot parse testserver URL port: %w", err)
		return
	}

	e = Endpoint{
		Host:    u.Hostname(),
		Port:    port,
		WithSSL: u.Scheme == "https",
	}

	return
}
