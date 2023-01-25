package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/stretchr/testify/assert"
)

func TestHeadBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	client := New(cfg)
	bucket := "myBucket"

	_, err = client.HeadBucket(context.Background(), bucket)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	s3client := s3.NewFromConfig(cfg.ToAWS())
	_, err = s3client.HeadBucket(
		context.Background(),
		&s3.HeadBucketInput{
			Bucket: &bucket,
		},
	)
	assert.NoError(t, err)
}
