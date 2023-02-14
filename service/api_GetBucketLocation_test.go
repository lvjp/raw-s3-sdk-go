package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/lvjp/raw-s3-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestGetBucketLocation(t *testing.T) {
	const input = `<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-west-1</LocationConstraint>`
	var expected = types.BucketLocationConstraint{LocationConstraint: "us-west-1"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := []byte(input)

		headers := w.Header()
		headers.Set("Content-Type", "application/xml")
		headers.Set("Content-Length", strconv.Itoa(len(payload)))

		w.WriteHeader(http.StatusOK)
		written, err := w.Write(payload)
		require.NoError(t, err)
		require.Equal(t, len(payload), written)
	}))
	defer ts.Close()

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

	client := New(cfg)

	output, err := client.GetBucketLocation(context.Background(), "pouet")
	require.NoError(t, err)
	require.Equal(t, expected, output.Payload)

	s3client := s3.NewFromConfig(cfg.ToAWS())

	s3out, err := s3client.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{Bucket: aws.String("pouet")})
	require.NoError(t, err)

	s3out.ResultMetadata = middleware.Metadata{}
	require.Equal(t, output.Payload.ToAWS(t), s3out.LocationConstraint)
}
