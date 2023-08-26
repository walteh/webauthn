all: binaries

build:
    ./hack/build

shell:
    ./hack/shell

binaries:
    docker buildx bake binaries

binaries-cross:
    docker buildx bake binaries-cross

# install: binaries
#     mkdir -p ~/bin
#     install bin/build/your-app ~/bin/your-app

release BIN_VERSION="local":
    BIN_VERSION={{BIN_VERSION}} ./hack/release

validate-all: lint test validate-vendor validate-docs validate-gen

lint:
    docker buildx bake lint

validate-vendor:
    docker buildx bake validate-vendor

validate-docs:
    docker buildx bake validate-docs

validate-gen:
    docker buildx bake validate-gen

update-all: vendor docs gen

vendor:
    ./hack/update-vendor

docs:
    ./hack/update-docs

outdated:
	docker buildx bake outdated
	cat ./bin/outdated/outdated.txt


gen:
    docker buildx bake update-gen --progress plain


test-driver:
    ./hack/test-driver
test:
    ./hack/test

test-unit:
    TESTPKGS=./... SKIP_INTEGRATION_TESTS=1 ./hack/test

test-integration:
    TESTPKGS=./tests ./hack/test


local:
	docker buildx bake image-default --progress plain


meta:
    docker buildx bake meta  --progress plain


install: binaries
	./bin/build/buildrc install && buildrc --version
