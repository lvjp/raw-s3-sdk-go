package service

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/lvjp/raw-s3-sdk-go/types"
)

type GetBucketLocationOutput struct {
	Payload types.BucketLocationConstraint

	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

func (c *Service) GetBucketLocation(ctx context.Context, bucket string) (*GetBucketLocationOutput, error) {
	req, resp, err := c.Do(ctx, http.MethodGet, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &GetBucketLocationOutput{
		HTTPRequest:  req,
		HTTPResponse: resp,
	}

	if err := xml.Unmarshal(body, &output.Payload); err != nil {
		return nil, err
	}

	return output, nil
}
