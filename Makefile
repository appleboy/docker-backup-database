DIST := dist
SERVICE ?= backup

DOCKER_ACCOUNT := appleboy
GOFMT ?= gofmt "-s"
SHASUM ?= shasum -a 256
GO ?= go
TARGETS ?= linux darwin windows
ARCHS ?= amd64
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOFILES := $(shell find . -name "*.go" ! -name "generated.*" ! -name "base.go" -type f)
TAGS ?= sqlite sqlite_unlock_notify json1

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
else
	VERSION ?= $(shell git describe --tags --always | sed 's/-/+/' | sed 's/^v//')
endif

LDFLAGS ?= -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)

all: build

.PHONY: generate
generate:
	@which fileb0x > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/UnnoTed/fileb0x; \
	fi
	$(GO) generate ./...

.PHONY: vendor
vendor:
	GO111MODULE=on $(GO) mod tidy && GO111MODULE=on $(GO) mod vendor

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

embedmd:
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/campoy/embedmd; \
	fi
	embedmd -d *.md

vet:
	$(GO) vet ./...

.PHONY: lint
lint:
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml ./... || exit 1

.PHONY: golangci-lint
golangci-lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		export BINARY="golangci-lint"; \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.18.0; \
	fi
	golangci-lint run --deadline=3m

install: $(GOFILES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(SERVICE)

$(SERVICE): $(GOFILES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@ ./cmd/$(SERVICE)

build_binary:
	$(GO) build -v -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/$(GOOS)/$(GOARCH)/$(DOCKER_IMAGE) ./cmd/$(SERVICE)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

upx:
	@hash upx > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Missing upx command"; \
		exit 1; \
	fi
	upx -o bin/$(SERVICE)-small bin/$(SERVICE)
	mv bin/$(SERVICE)-small bin/$(SERVICE)

.PHONY: unit-test-coverage
unit-test-coverage:
	@$(GO) test -v -race -cover -coverprofile coverage.out -tags '$(TAGS)' ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

test:
	@$(GO) test -cover -tags '$(TAGS)' ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

release: release-dirs release-build release-copy release-compress release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	gox -os="$(TARGETS)" -arch="$(ARCHS)" -tags="$(TAGS)" -ldflags="$(EXTLDFLAGS)-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(SERVICE)-$(VERSION)-{{.OS}}-{{.Arch}}" ./cmd/$(SERVICE)/...

.PHONY: release-compress
release-compress:
	@hash gxz > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/ulikunitz/xz/cmd/gxz; \
	fi
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && gxz -k -9 $${file}; done;

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(SERVICE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;

clean_dist:
	rm -rf bin dist

clean: clean_dist
	$(GO) clean -modcache -cache -x -i ./...
