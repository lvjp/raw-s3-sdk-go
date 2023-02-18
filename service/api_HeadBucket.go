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

	req, res, err := s.doCall(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	output.HTTPRequest = req
	output.HTTPResponse = res

	return &output, nil
}
