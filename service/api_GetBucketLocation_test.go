package service

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/lvjp/raw-s3-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestGetBucketLocation(t *testing.T) {
	var expected = types.LocationConstraint{LocationConstraint: "us-west-1"}

	handler := NewSimpleXMLResponseHandler(t, &expected)
	ts, ourClient, awsClient := NewServer(t, handler)
	defer ts.Close()

	bucket := "myBucket"

	t.Run("our", func(t *testing.T) {
		ourOutput, err := ourClient.GetBucketLocation(context.Background(), bucket)
		require.NoError(t, err)
		require.Equal(t, expected, ourOutput.Payload)
	})

	t.Run("aws", func(t *testing.T) {
		s3out, err := awsClient.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{Bucket: &bucket})
		require.NoError(t, err)

		s3out.ResultMetadata = middleware.Metadata{}
		require.Equal(t, expected.ToAWS(t), &s3out.LocationConstraint)
	})
}
