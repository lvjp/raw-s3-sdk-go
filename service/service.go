package service

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/lvjp/raw-s3-sdk-go/signing"
)

const standardHTTPPort = 80
const standardHTTPSPort = 443

type Service struct {
	config config.Config
}

func New(config config.Config) *Service {
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}

	return &Service{
		config: config,
	}
}

func (s *Service) Do(ctx context.Context, method string, bucket, key *string, queryString url.Values, body io.ReadCloser) (*http.Request, *http.Response, error) {
	req := s.newRequest(ctx, method, bucket, key, queryString, body)

	signer, err := signing.NewSigner(s.config.SignatureType)
	if err != nil {
		return nil, nil, err
	}

	err = signer(req, s.config.Credentials, s.config.Region)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.config.HTTPClient.Do(req)

	return req, resp, err
}

func (s *Service) newRequest(ctx context.Context, method string, bucket, key *string, queryString url.Values, body io.ReadCloser) *http.Request {
	url := s.newURL(bucket, key, queryString)

	req := &http.Request{
		Method:     method,
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"User-Agent": []string{"raw-s3-sdk-go"},
		},
		Body: body,
		Host: url.Host,
	}

	return req.WithContext(ctx)
}

func (s *Service) newURL(bucket *string, key *string, queryString url.Values) *url.URL {
	e := &s.config.Endpoint

	url := &url.URL{
		Host: e.Host,
	}

	if e.WithVirtualHost && bucket != nil {
		url.Host = *bucket + url.Host
	}

	if e.WithSSL {
		url.Scheme = "https"
		if e.Port != standardHTTPSPort {
			url.Host += ":" + strconv.Itoa(e.Port)
		}
	} else {
		url.Scheme = "http"
		if e.Port != standardHTTPPort {
			url.Host += ":" + strconv.Itoa(e.Port)
		}
	}

	return url
}

func (s *Service) doCall(ctx context.Context, respBody any) (*http.Request, *http.Response, error) {
	req, resp, err := s.Do(ctx, http.MethodGet, nil, nil, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if respBody != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		if err := xml.Unmarshal(body, respBody); err != nil {
			return nil, nil, err
		}
	}

	return req, resp, nil
}
