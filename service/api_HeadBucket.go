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

func (c *Service) HeadBucket(ctx context.Context, bucket string) (*HeadBucketOutput, error) {
	req, resp, err := c.Do(ctx, http.MethodHead, &bucket, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	output := &HeadBucketOutput{
		HTTPRequest:  req,
		HTTPResponse: resp,
	}

	return output, nil
}
