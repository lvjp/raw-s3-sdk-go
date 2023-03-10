package types

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"
)

var _ AWSConvertible[s3.ListBucketsOutput] = (*ListAllMyBucketsResult)(nil)

type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
	Owner   *Owner
}

type Bucket struct {
	CreationDate *string
	Name         *string
}

type Owner struct {
	DisplayName *string
	ID          *string
}

func (lambr *ListAllMyBucketsResult) ToAWS(t *testing.T) *s3.ListBucketsOutput {
	result := &s3.ListBucketsOutput{}

	if lambr.Buckets != nil {
		result.Buckets = make([]types.Bucket, 0, len(lambr.Buckets))
		for _, bucket := range lambr.Buckets {
			result.Buckets = append(result.Buckets, *bucket.ToAWS(t))
		}
	}

	if lambr.Owner != nil {
		result.Owner = lambr.Owner.ToAWS()
	}

	return result
}

func (b *Bucket) ToAWS(t *testing.T) *types.Bucket {
	result := &types.Bucket{
		Name: b.Name,
	}

	if b.CreationDate != nil {
		parsed, err := time.Parse(time.RFC3339, *b.CreationDate)
		require.NoError(t, err, "Cannot parse bucket '%v' creation date", b.Name)
		result.CreationDate = &parsed
	}

	return result
}

func (o *Owner) ToAWS() *types.Owner {
	return &types.Owner{
		DisplayName: o.DisplayName,
		ID:          o.ID,
	}
}
