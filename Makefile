DIST := dist
SERVICE ?= backup
GOFMT ?= gofumpt -l -s
SHASUM ?= shasum -a 256
GO ?= go
TARGETS ?= linux darwin windows
ARCHS ?= amd64
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOFILES := $(shell find . -name "*.go" -type f)
TAGS ?=
GOPATH ?= $(shell $(GO) env GOPATH)

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

LDFLAGS ?= -X cmd/$(SERVICE)/main.Version=$(VERSION) -X cmd/$(SERVICE)/main.BuildDate=$(BUILD_DATE)

all: build

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: vendor
vendor:
	GO111MODULE=on $(GO) mod tidy && GO111MODULE=on $(GO) mod vendor

.PHONY: fmt
fmt:
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u mvdan.cc/gofumpt; \
	fi
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u mvdan.cc/gofumpt; \
	fi
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

.PHONY: lint
lint:
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1; \
	fi
	golangci-lint run -v --deadline=3m

install: $(GOFILES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(SERVICE)

$(SERVICE): $(GOFILES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@ ./cmd/$(SERVICE)

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

build_binary:
	$(GO) build -v -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/$(GOOS)/$(GOARCH)/$(SERVICE) ./cmd/$(SERVICE)

clean_dist:
	rm -rf bin release

clean: clean_dist
	$(GO) clean -modcache -cache -x -i ./...
