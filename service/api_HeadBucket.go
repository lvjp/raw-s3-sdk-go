package service

import (
	"context"
	"net/http"
)

type HeadBucketInput struct {
	Bucket string
}

type HeadBucketOutput struct {
	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

func (s *Service) HeadBucket(ctx context.Context, bucket string) (*HeadBucketOutput, error) {
	output := HeadBucketOutput{}
	var err error

	output.HTTPRequest, output.HTTPResponse, err = s.doCall(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
