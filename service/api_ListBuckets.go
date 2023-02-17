package service

import (
	"context"
	"net/http"

	"github.com/lvjp/raw-s3-sdk-go/types"
)

type ListBucketsOutput struct {
	Payload types.ListAllMyBucketsResult

	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

func (s *Service) ListBuckets(ctx context.Context) (*ListBucketsOutput, error) {
	output := ListBucketsOutput{}
	var err error

	output.HTTPRequest, output.HTTPResponse, err = s.doCall(ctx, &output.Payload)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
