package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
			CreationDate: "2019-12-11T23:32:47+00:00",
			Name:         "DOC-EXAMPLE-BUCKET",
		},
		{
			CreationDate: "2019-11-10T23:32:13+00:00",
			Name:         "DOC-EXAMPLE-BUCKET2",
		},
	},
	Owner: types.Owner{
		DisplayName: "Account+Name",
		ID:          "DUMMYACKCEVSQ6C2EXAMPLE",
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
}
