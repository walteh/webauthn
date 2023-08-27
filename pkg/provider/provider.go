package provider

import (
	"context"
	"reflect"

	"github.com/walteh/webauthn/pkg/webauthn/types"
)

type ManagedObject string

const (
	Challenge  ManagedObject = "ceremony"
	Credential ManagedObject = "credential"
)

type Provider interface {
	Get(ctx context.Context, challenge string, cred string) (*types.Credential, error)
}

type Validation struct {
	Type  ManagedObject
	ID    string
	Field reflect.StructField
	Value interface{}
}

type Update struct{}
