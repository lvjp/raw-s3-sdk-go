package signing

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/lvjp/raw-s3-sdk-go/session"
	"github.com/stretchr/testify/assert"
)

// github.com/stretchr/testify v1.8.1
// golang.org/x/exp v0.0.0-20230116083435-1de6713980de

func TestSign(t *testing.T) {
	for i, tc := range generateTestCases(t) {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := Sign(
				tc.request,
				session.Credentials{
					AccessKey: "AKIAIOSFODNN7EXAMPLE",
					SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
				tc.region,
			)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.authorizationHeader, tc.request.Header.Get("Authorization"))
			}
		})
	}
}

func TestSign_insider(t *testing.T) {
	for i, tc := range generateTestCases(t) {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prepareRequest(tc.request)

			signer := signer{
				request: tc.request,
				credentials: &session.Credentials{
					AccessKey: "AKIAIOSFODNN7EXAMPLE",
					SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				},
				region:      tc.region,
				queryString: tc.request.URL.Query(),
			}
			t.Run("Date", func(t *testing.T) {
				if assert.NoError(t, signer.extractDate()) {
					assert.Equal(t, tc.date, signer.date)
				}
			})

			t.Run("Scope", func(t *testing.T) {
				signer.computeScope()
				assert.Equal(t, tc.scope, signer.scope)
			})

			t.Run("CanonicalRequest", func(t *testing.T) {
				// t.Skip()
				assert.Equal(t, tc.canonicalRequest, signer.computeCanonicalRequest())
			})

			var stringToSign string
			t.Run("StringToSign", func(t *testing.T) {
				// t.Skip()
				stringToSign = signer.computeStringToSign()
				assert.Equal(t, tc.stringToSign, stringToSign)
			})

			var signature string
			t.Run("Signature", func(t *testing.T) {
				// t.Skip()
				signature = signer.computeSignature(signer.computeSigningKey(), stringToSign)
				assert.Equal(t, tc.signature, signature)
			})

			t.Run("Authorization", func(t *testing.T) {
				header := signer.computeAuthorizationheader(signature)
				assert.Equal(t, tc.authorizationHeader, header)
			})
		})
	}
}

func newRequest(t *testing.T, method, url string, headers map[string]string) *http.Request {
	r, err := http.NewRequest(method, url, nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	return r
}

type testCase struct {
	request *http.Request
	region  string

	date                time.Time
	scope               string
	canonicalRequest    string
	stringToSign        string
	signature           string
	authorizationHeader string
}

func generateTestCases(t *testing.T) []testCase {
	return []testCase{
		// Example: GET Object
		{
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
			region: "us-east-1",
			date:   time.Date(2013, time.May, 24, 0, 0, 0, 0, time.UTC),
			scope:  "20130524/us-east-1/s3/aws4_request",
			canonicalRequest: `GET
/test.txt

host:examplebucket.s3.amazonaws.com
range:bytes=0-9
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;range;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
			stringToSign: `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
7344ae5b7ee6c3e7e6b0fe0640412a37625d1fbfff95c48bbb2dc43964946972`,
			signature:           "f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41",
			authorizationHeader: "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;range;x-amz-content-sha256;x-amz-date,Signature=f0e8bdb87c964420e857bd35b5d6ed310bd44f0170aba48dd91039c6036bdb41",
		},

		// Example: PUT Object
		{
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
			region: "us-east-1",
			date:   time.Date(2013, time.May, 24, 0, 0, 0, 0, time.UTC),
			scope:  "20130524/us-east-1/s3/aws4_request",
			canonicalRequest: `PUT
/test%24file.text

date:Fri, 24 May 2013 00:00:00 GMT
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072
x-amz-date:20130524T000000Z
x-amz-storage-class:REDUCED_REDUNDANCY

date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class
44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072`,
			stringToSign: `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9e0e90d9c76de8fa5b200d8c849cd5b8dc7a3be3951ddb7f6a76b4158342019d`,
			signature:           "98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd",
			authorizationHeader: "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request,SignedHeaders=date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class,Signature=98ad721746da40c64f1a55b78f14c238d841ea1380cd77a1b5971af0ece108bd",
		},

		// Example: GET Bucket Lifecycle
		{
			request: newRequest(
				t,
				http.MethodGet,
				"http://examplebucket.s3.amazonaws.com?lifecycle",
				map[string]string{
					"x-amz-date":           "20130524T000000Z",
					"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			),
			region: "us-east-1",
			date:   time.Date(2013, time.May, 24, 0, 0, 0, 0, time.UTC),
			scope:  "20130524/us-east-1/s3/aws4_request",
			canonicalRequest: `GET
/
lifecycle=
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
			stringToSign: `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9766c798316ff2757b517bc739a67f6213b4ab36dd5da2f94eaebf79c77395ca`,
			signature:           "fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543",
			authorizationHeader: "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;x-amz-content-sha256;x-amz-date,Signature=fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543",
		},

		// Example: Get Bucket (List Objects)
		{
			request: newRequest(
				t,
				http.MethodGet,
				"http://examplebucket.s3.amazonaws.com?max-keys=2&prefix=J",
				map[string]string{
					"x-amz-date":           "20130524T000000Z",
					"x-amz-content-sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				},
			),
			region: "us-east-1",
			date:   time.Date(2013, time.May, 24, 0, 0, 0, 0, time.UTC),
			scope:  "20130524/us-east-1/s3/aws4_request",
			canonicalRequest: `GET
/
max-keys=2&prefix=J
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`,
			stringToSign: `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
df57d21db20da04d7fa30298dd4488ba3a2b47ca3a489c74750e0f1e7df1b9b7`,
			signature:           "34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7",
			authorizationHeader: "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request,SignedHeaders=host;x-amz-content-sha256;x-amz-date,Signature=34b48302e7b5fa45bde8084f4b7868a86f0a534bc59db6670ed5711ef69dc6f7",
		},
	}
}
