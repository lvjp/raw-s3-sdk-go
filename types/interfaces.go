package types

import (
	"testing"
)

type AWSConvertible[T any] interface {
	ToAWS(t *testing.T) *T
}
