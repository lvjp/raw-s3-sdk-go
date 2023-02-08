package signing

import (
	"net/http"
	"testing"

	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/stretchr/testify/require"
)

var creds = config.Credentials{
	AccessKey: "AKI" + "AIOSFODNN7EXAMPLE",
	SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
}

func TestSign(t *testing.T) {
	for _, tc := range generateTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			err := Sign(
				tc.request,
				creds,
				tc.region,
			)

			require.NoError(t, err)
			require.Equal(t, tc.expected, tc.request.Header.Get("Authorization"))
		})
	}
}

func newRequest(t *testing.T, method, url string, headers map[string]string) *http.Request {
	r, err := http.NewRequest(method, url, nil)
	require.NoError(t, err)

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	return r
}

type testCase struct {
	name     string
	request  *http.Request
	region   string
	expected string
}

func generateTestCases(t *testing.T) []testCase {
	return []testCase{
		{
			name: "GetObject",
			request: newRequest(
				t,
				http.MethodGet,
				"http://examplebucket.s3.amazonaws.com/test.txt",
				map[string]string{
					"Range":                "bytes=0-9",
					"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					"x-amz-date":           "20130524T000000Z",
				},
			),
			region:   "us-east-1",
			expected: "AWS4-HMAC-SHA256 Credential=" + creds.AccessKey + "/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;range;x-amz-content-sha256;x-amz-date,Signature=f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41",
		},

		{
			name: "PutObject",
			request: newRequest(
				t,
				http.MethodPut,
				"http://examplebucket.s3.amazonaws.com/test$file.text",
				map[string]string{
					"Date":                 "Fri, 24 May 2013 00:00:00 GMT",
					"x-amz-date":           "20130524T000000Z",
					"x-amz-storage-class":  "REDUCED_REDUNDANCY",
					"x-amz-content-sha256": "44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072",
				},
			),
			region:   "us-east-1",
			expected: "AWS4-HMAC-SHA256 Credential=" + creds.AccessKey + "/20130524/us-east-1/s3/aws4_request,SignedHeaders=date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class,Signature=98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd",
		},

		{
			name: "GetBucketLifecycle",
			request: newRequest(
				t,
				http.MethodGet,
				"http://examplebucket.s3.amazonaws.com?lifecycle",
				map[string]string{
					"x-amz-date":           "20130524T000000Z",
					"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			),
			region:   "us-east-1",
			expected: "AWS4-HMAC-SHA256 Credential=" + creds.AccessKey + "/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;x-amz-content-sha256;x-amz-date,Signature=fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543",
		},

		{
			name: "ListObjects",
			request: newRequest(
				t,
				http.MethodGet,
				"http://examplebucket.s3.amazonaws.com?max-keys=2&prefix=J",
				map[string]string{
					"x-amz-date":           "20130524T000000Z",
					"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			),
			region:   "us-east-1",
			expected: "AWS4-HMAC-SHA256 Credential=" + creds.AccessKey + "/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;x-amz-content-sha256;x-amz-date,Signature=34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7",
		},
	}
}
