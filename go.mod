module github.com/walteh/webauthn

go 1.21.5

replace github.com/walteh/snake/szerolog => ../snake/szerolog

require (
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/credentials v1.13.37
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.16.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.21.5
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/go-webauthn/webauthn v0.9.4
	github.com/go-webauthn/x v0.1.5
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.4.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.31.0
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.8.4
	github.com/ugorji/go/codec v1.2.11
	github.com/walteh/terrors v0.10.0
	golang.org/x/crypto v0.16.0
	golang.org/x/net v0.15.0
	golang.org/x/sync v0.5.0
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
