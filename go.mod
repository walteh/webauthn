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
	github.com/walteh/buildrc v0.18.1
	github.com/walteh/snake v0.40.0
	github.com/walteh/terrors v0.10.0
	golang.org/x/crypto v0.16.0
	golang.org/x/net v0.15.0
	golang.org/x/sync v0.5.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.39 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.39 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.42 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.15.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.15.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.21.5 // indirect
	github.com/aws/smithy-go v1.14.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/containerd/continuity v0.4.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/cli v24.0.5+incompatible // indirect
	github.com/docker/docker v24.0.5-0.20230714235725-36e9e796c6fc+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.4.7 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/opencontainers/runc v1.1.9 // indirect
	github.com/ory/dockertest/v3 v3.10.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.0 // indirect
)
