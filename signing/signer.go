package signing

import (
	"fmt"
	"net/http"

	"github.com/lvjp/raw-s3-sdk-go/config"
	signv4 "github.com/lvjp/raw-s3-sdk-go/signing/v4"
)

type Signer func(r *http.Request, credentials config.Credentials, region string) error

func NewSigner(signatureType config.SignatureType) (Signer, error) {
	switch signatureType {
	case config.SignatureTypeV4:
		return signv4.Sign, nil
	default:
		return nil, fmt.Errorf("unsupported signature type: %v", signatureType)
	}
}
