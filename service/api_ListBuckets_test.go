package service

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/lvjp/raw-s3-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func TestListBuckets(t *testing.T) {
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

	handler := NewSimpleXMLResponseHandler(t, &expected)
	ts, ourClient, awsClient := NewServer(t, handler)
	defer ts.Close()

	t.Run("our", func(t *testing.T) {
		output, err := ourClient.ListBuckets(context.Background())
		require.NoError(t, err)
		require.Equal(t, expected, output.Payload)
	})

	t.Run("aws", func(t *testing.T) {
		s3out, err := awsClient.ListBuckets(context.Background(), &s3.ListBucketsInput{})
		require.NoError(t, err)

		s3out.ResultMetadata = middleware.Metadata{}
		require.Equal(t, expected.ToAWS(t), s3out)
	})
}
