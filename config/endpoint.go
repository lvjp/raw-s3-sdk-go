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

func (e Endpoint) String() string {
	suffix := ""

	if e.WithSSL {
		suffix = "s"
	}

	return fmt.Sprintf("http%s://%s:%d", suffix, e.Host, e.Port)
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
