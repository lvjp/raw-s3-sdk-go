package session

type SignatureType int

const (
	SignatureTypeV2       SignatureType = 1
	SignatureTypeV2Header SignatureType = 2
	SignatureTypeV4       SignatureType = 3
)

type Config struct {
	Region string

	Endpoint struct {
		Host string
		Port int

		WithSSL         bool
		WithVirtualHost bool
	}

	Credentials   Credentials
	SignatureType SignatureType
}

type Credentials struct {
	AccessKey  string
	SecretKey string
}
