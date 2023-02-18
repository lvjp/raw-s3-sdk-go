package service

import (
	"context"
	"net/http"

	"github.com/lvjp/raw-s3-sdk-go/types"
)

type GetBucketLocationOutput struct {
	Payload types.LocationConstraint

	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

func (s *Service) GetBucketLocation(ctx context.Context, bucket string) (*GetBucketLocationOutput, error) {
	output := GetBucketLocationOutput{}

	req, res, err := s.doCall(ctx, &output.Payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	output.HTTPRequest = req
	output.HTTPResponse = res

	return &output, nil
}
