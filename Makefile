NAME=forward-server
BUILDDIR=build

BASEPATH := $(shell pwd)
BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
BUILD := $(shell git rev-parse --short HEAD)
VERSION ?= $(BRANCH)-$(BUILD)
BuildTime:= $(shell date -u '+%Y-%m-%d %I:%M:%S%p')
COMMIT:= $(shell git rev-parse HEAD)
GOVERSION:= $(shell go version)

LDFLAGS+=-X 'main.BuildStamp=$(BuildTime)'
LDFLAGS+=-X 'main.GitHash=$(COMMIT)'
LDFLAGS+=-X 'main.GoVersion=$(GOVERSION)'
LDFLAGS+=-X 'main.Version=$(VERSION)'

GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)"

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm64

all-arch: $(PLATFORM_LIST)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@ .
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@ .
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@ .
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BUILDDIR)/$(NAME)-$@ .
	mkdir -p $(BUILDDIR)/$(NAME)-$(VERSION)-$@/
	cp $(BUILDDIR)/$(NAME)-$@ $(BUILDDIR)/$(NAME)-$(VERSION)-$@/$(NAME)
	rm -rf $(BUILDDIR)/$(NAME)-$(VERSION)-$@ $(BUILDDIR)/$(NAME)-$@

.PHONY: clean
clean:
	-rm -rf $(BUILDDIR)
