package service

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/lvjp/raw-s3-sdk-go/types"
)

type ListBucketsOutput struct {
	Payload types.ListAllMyBucketsResult

	HttpRequest  *http.Request
	HttpResponse *http.Response
}

func (c *Service) ListBuckets(ctx context.Context) (*ListBucketsOutput, error) {
	req, resp, err := c.Do(ctx, http.MethodGet, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := &ListBucketsOutput{
		HttpRequest:  req,
		HttpResponse: resp,
	}

	if err := xml.Unmarshal(body, &output.Payload); err != nil {
		return nil, err
	}

	return output, nil
}
