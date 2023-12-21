
.PHONY: generate lint hash

vendor_sha := $(shell find ./tools/vendor -type f | xargs sha256sum -b | sha256sum | cut -d ' ' -f1 | head)

vendor_sha := $(shell find ./tools/vendor -type f | xargs sha256sum -b | sha256sum | cut -d ' ' -f1 | head)


./tools/bin/$(vendor_sha).make-hash:
	@echo $(vendor_sha)
	@mkdir -p ./tools/bin
	@echo $(vendor_sha) > ./tools/bin/$(vendor_sha).make-hash

./tools/bin/buf: ./tools/bin/$(vendor_sha).make-hash
	cd tools && GOWORK=off go build -mod=vendor -o ./bin/buf github.com/bufbuild/buf/cmd/buf

./tools/bin/mockery: ./tools/bin/$(vendor_sha).make-hash
	cd tools && GOWORK=off go build -mod=vendor -o ./bin/mockery github.com/vektra/mockery/v2

./tools/bin/golangci-lint: ./tools/bin/$(vendor_sha).make-hash
	cd tools && GOWORK=off go build -mod=vendor -o ./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

./gen/buf/$(vendor_sha).make-hash:  ./tools/bin/buf
		./tools/bin/buf generate --include-imports --include-wkt --exclude-path="tools/vendor" --template ./buf.gen.yaml --path=./proto  && \
			echo $(vendor_sha) > ./gen/buf/$(vendor_sha).make-hash

./gen/mockery/$(vendor_sha).make-hash:  ./tools/bin/mockery
		./tools/bin/mockery --dir ./gen/mockery && \
			echo $(vendor_sha) > ./gen/mockery/$(vendor_sha).make-hash

lint:  ./tools/bin/golangci-lint
		GOWORK=off ./tools/bin/golangci-lint run --config ./.golangci.yml

generate: ./gen/buf/$(vendor_sha).make-hash
generate: ./gen/mockery/$(vendor_sha).make-hash
