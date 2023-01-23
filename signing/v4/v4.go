package signing

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/lvjp/raw-s3-sdk-go/config"
	"github.com/lvjp/raw-s3-sdk-go/signing/utils"
	"golang.org/x/exp/maps"
)

const dateFormatYYYMMDD = "20060102"
const dateFormatISO8601 = "20060102T150405Z"

func Sign(r *http.Request, credentials config.Credentials, region string) error {
	prepareRequest(r)

	signer := &signer{
		request:     r,
		credentials: &credentials,
		region:      region,
		queryString: r.URL.Query(),
	}

	if err := signer.extractDate(); err != nil {
		return err
	}

	signer.computeScope()

	auth := signer.computeAuthorizationheader(
		signer.computeSignature(
			signer.computeSigningKey(),
			signer.computeStringToSign(),
		),
	)

	r.Header.Set("Authorization", auth)

	return nil
}

func prepareRequest(r *http.Request) {
	defaults := map[string]string{
		"Host":       r.Host,
		"X-Amz-Date": time.Now().UTC().Format(dateFormatISO8601),
	}

	for name, value := range defaults {
		if r.Header.Get(name) == "" {
			r.Header.Set(name, value)
		}
	}

	if r.Header.Get("X-Amz-Content-Sha256") == "" {
		payload := readAndReplaceBody(r)
		hash := fmt.Sprintf("%x", sha256.Sum256(payload))
		r.Header.Set("X-Amz-Content-Sha256", hash)
	}

	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
}

type signer struct {
	// Input
	request     *http.Request
	credentials *config.Credentials
	region      string
	queryString url.Values

	// Cached computed values
	date          time.Time
	scope         string
	signedHeaders string
}

func (s *signer) extractDate() (err error) {
	for _, header := range []string{"x-amz-date", "date"} {
		if raw := s.request.Header.Get(header); raw != "" {
			s.date, err = time.Parse(dateFormatISO8601, raw)
			if err != nil {
				err = fmt.Errorf("cannot parse the header '%s: %s': %w", header, raw, err)
			}
			return
		}
	}

	if raw := s.queryString.Get("x-amz-date"); raw != "" {
		s.date, err = time.Parse(dateFormatISO8601, raw)
		if err != nil {
			err = fmt.Errorf("cannot parse the header query parameter 'x-amz-date': %w", err)
			return
		}
	}

	return errors.New("cannot find date for the signature")
}

func (s *signer) computeScope() {
	s.scope = fmt.Sprintf(
		"%s/%s/s3/aws4_request",
		s.date.Format(dateFormatYYYMMDD),
		s.region,
	)
}

func (s *signer) computeSigningKey() []byte {
	res := []byte("AWS4" + s.credentials.SecretKey)

	toSign := []string{s.date.Format(dateFormatYYYMMDD), s.region, "s3", "aws4_request"}
	for _, data := range toSign {
		res = utils.HMacSha256(res, data)
	}

	return res
}

func (s *signer) computeCanonicalRequest() string {
	return strings.Join(
		[]string{
			s.request.Method,
			utils.UriEncode(s.request.URL.Path),
			s.computeCanonicalQueryString(),
			s.computeCanonicalHeaders(),
			s.request.Header.Get("X-Amz-Content-Sha256"),
		},
		"\n",
	)
}

func (s *signer) computeStringToSign() string {
	return strings.Join(
		[]string{
			"AWS4-HMAC-SHA256",
			s.date.Format(dateFormatISO8601),
			s.scope,
			utils.HexSha256(s.computeCanonicalRequest()),
		},
		"\n",
	)
}

func (s *signer) computeSignature(signingKey []byte, stringToSign string) string {
	return hex.EncodeToString(utils.HMacSha256(signingKey, stringToSign))
}

func (s *signer) computeAuthorizationheader(signature string) string {
	return fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s,SignedHeaders=%s,Signature=%s",
		s.credentials.AccessKey,
		s.scope,
		s.signedHeaders,
		signature,
	)
}

func (s *signer) computeCanonicalQueryString() string {
	keys := maps.Keys(s.queryString)
	sort.Strings(keys)

	buf := strings.Builder{}
	for _, key := range keys {
		values := s.queryString[key]
		keyEscaped := utils.UriEncode(key)

		for _, value := range values {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(utils.UriEncode(value))
		}
	}

	return buf.String()
}

func (s *signer) computeCanonicalHeaders() string {
	headers := make(map[string][]string, len(s.request.Header))

	for name, values := range s.request.Header {
		name := strings.ToLower(name)
		switch name {
		case "host":
			cleaned := make([]string, 0, len(values))
			for _, value := range values {
				if strings.Contains(value, ":") {
					split := strings.Split(value, ":")
					port := split[1]
					if port == "80" || port == "443" {
						value = split[0]
					}
				}
				cleaned = append(cleaned, value)
			}
			headers[name] = cleaned
		case "content-type", "date", "range":
			headers[name] = values
		default:
			if strings.HasPrefix(name, "x-amz-") {
				headers[name] = values
			}
		}
	}

	keys := maps.Keys(headers)
	sort.Strings(keys)
	s.signedHeaders = strings.Join(keys, ";")

	buf := strings.Builder{}
	for _, key := range keys {
		values := headers[key]
		keyEscaped := utils.UriEncode(key)

		for _, value := range values {
			buf.WriteString(keyEscaped)
			buf.WriteByte(':')
			buf.WriteString(strings.TrimSpace(value))
			buf.WriteByte('\n')
		}
	}

	buf.WriteByte('\n')
	buf.WriteString(s.signedHeaders)

	return buf.String()
}

func readAndReplaceBody(r *http.Request) []byte {
	if r.Body == nil {
		return []byte{}
	}

	payload, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(payload))
	return payload
}
