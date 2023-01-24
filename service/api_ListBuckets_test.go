package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/lvjp/raw-s3-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

const input = `<ListAllMyBucketsResult>
<Buckets>
  <Bucket>
	<CreationDate>2019-12-11T23:32:47+00:00</CreationDate>
	<Name>DOC-EXAMPLE-BUCKET</Name>
  </Bucket>
  <Bucket>
	<CreationDate>2019-11-10T23:32:13+00:00</CreationDate>
	<Name>DOC-EXAMPLE-BUCKET2</Name>
  </Bucket>
</Buckets>
<Owner>
  <DisplayName>Account+Name</DisplayName>
  <ID>DUMMYACKCEVSQ6C2EXAMPLE</ID>
</Owner>
</ListAllMyBucketsResult>`

var expected = types.ListAllMyBucketsResult{
	Buckets: []types.Bucket{
		{
			CreationDate: aws.String("2019-12-11T23:32:47+00:00"),
			Name:         aws.String("DOC-EXAMPLE-BUCKET"),
		},
		{
			CreationDate: aws.String("2019-11-10T23:32:13+00:00"),
			Name:         aws.String("DOC-EXAMPLE-BUCKET2"),
		},
	},
	Owner: &types.Owner{
		DisplayName: aws.String("Account+Name"),
		ID:          aws.String("DUMMYACKCEVSQ6C2EXAMPLE"),
	},
}

func TestListBucket(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, input)
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

	output, err := client.ListBuckets(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expected, output.Payload)
	}

	s3client := s3.NewFromConfig(cfg.ToAWS())

	s3out, err := s3client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if assert.NoError(t, err) {
		s3out.ResultMetadata = middleware.Metadata{}
		assert.Equal(t, output.Payload.ToAWS(t), s3out)
	}
}
