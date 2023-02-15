package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/require"
)

func TestHeadBucket(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts, ourClient, awsClient := NewServer(t, handler)
	defer ts.Close()

	bucket := "myBucket"

	t.Run("our", func(t *testing.T) {
		_, err := ourClient.HeadBucket(context.Background(), bucket)
		require.NoError(t, err)
	})

	t.Run("aws", func(t *testing.T) {
		_, err := awsClient.HeadBucket(
			context.Background(),
			&s3.HeadBucketInput{
				Bucket: &bucket,
			},
		)
		require.NoError(t, err)
	})
}
