package types

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var _ AWSConvertible[types.BucketLocationConstraint] = (*LocationConstraint)(nil)

type LocationConstraint struct {
	LocationConstraint string `xml:",chardata"`
}

func (lc *LocationConstraint) ToAWS(t *testing.T) *types.BucketLocationConstraint {
	ret := types.BucketLocationConstraint(lc.LocationConstraint)
	return &ret
}
