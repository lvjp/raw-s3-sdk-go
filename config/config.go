package config

import "net/http"

type SignatureType int

const (
	SignatureTypeAnonymous SignatureType = 1
	SignatureTypeV2        SignatureType = 2
	SignatureTypeV2Header  SignatureType = 3
	SignatureTypeV4        SignatureType = 4
)

type Config struct {
	HttpClient HttpClient

	Region string

	Endpoint Endpoint

	Credentials   Credentials
	SignatureType SignatureType
}

type Credentials struct {
	AccessKey string
	SecretKey string
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}
