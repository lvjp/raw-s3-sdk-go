package service

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/stretchr/testify/require"
)

func NewServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Service, *s3.Client) {
	ts := httptest.NewServer(handler)

	cfg := config.Config{
		HTTPClient: ts.Client(),

		Region: "eu-west-1",

		Credentials: config.Credentials{
			AccessKey: "DUMMYAIOSFODNN7EXAMPLE",
			SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},

		SignatureType: config.SignatureTypeV4,
	}

	var err error
	cfg.Endpoint, err = config.NewEndpointFromURL(ts.URL)
	require.NoError(t, err)

	return ts, New(cfg), s3.NewFromConfig(cfg.ToAWS())
}

func NewSimpleXMLResponseHandler(t *testing.T, payload any) http.HandlerFunc {
	raw, err := xml.Marshal(payload)
	require.NoError(t, err, "Marshal payload")

	contentLength := strconv.Itoa(len(raw))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Content-Type", "application/xml")
		headers.Set("Content-Length", contentLength)

		w.WriteHeader(http.StatusOK)
		written, err := w.Write(raw)
		require.NoError(t, err)
		require.Equal(t, len(raw), written)
	})
}
