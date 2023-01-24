package config

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SignatureType int

const (
	SignatureTypeAnonymous SignatureType = 1
	SignatureTypeV2        SignatureType = 2
	SignatureTypeV2Header  SignatureType = 3
	SignatureTypeV4        SignatureType = 4
)

type Config struct {
	HTTPClient HTTPClient

	Region string

	Endpoint Endpoint

	Credentials   Credentials
	SignatureType SignatureType
}

func (c Config) ToAWS() aws.Config {
	return aws.Config{
		Region: c.Region,
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service != s3.ServiceID {
				return aws.Endpoint{}, errors.New("unsupported service: " + service)
			}

			return aws.Endpoint{
				URL:               c.Endpoint.String(),
				SigningRegion:     c.Region,
				HostnameImmutable: !c.Endpoint.WithVirtualHost,
			}, nil
		}),
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     c.Credentials.AccessKey,
				SecretAccessKey: c.Credentials.SecretKey,
			}, nil
		}),
		HTTPClient: c.HTTPClient,
	}
}

type Credentials struct {
	AccessKey string
	SecretKey string
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
