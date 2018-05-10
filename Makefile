#
GODEP := godep
GO := $(GODEP) go
GOTEST_FLAGS := -v
TEST_TIMEOUT := 180s

GOFILES := find . -name '*.go' ! -path './Godeps/*' ! -path './vendor/*'
GOFOLDERS := $(GO) list ./... | sed 's:^github.com/aristanetworks/goeapi:.:' | grep -vw -e './vendor' -e './examples'

# Code Coverage Related
COVER_TMPFILE := coverage.out
COVER_MODE := count

# External Tools
EXTERNAL_TOOLS=\
	golang.org/x/tools/cmd/cover

PKGS := $(shell go list ./... | grep -v /examples)

GOLINT := golint

all: install

install:
	$(GO) install $(PKGS)

test: unittest

systest:
	$(GOFOLDERS) | xargs $(GO) test $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) -run SystemTest$

unittest:
	$(GOFOLDERS) | xargs $(GO) test $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) -run UnitTest$

COVER_PKGS := `find . -name '*_test.go' ! -path "./.git/*" ! -path "./Godeps/*" ! -path "./vendor/*" | xargs -I{} dirname {} | sort -u`
coverdata:
	echo 'mode: $(COVER_MODE)' > $(COVER_TMPFILE)
	for dir in $(COVER_PKGS); do \
		$(GO) test -covermode=$(COVER_MODE) -coverprofile=$(COVER_TMPFILE).tmp -run UnitTest$  $$dir || exit; \
		tail -n +2 $(COVER_TMPFILE).tmp >> $(COVER_TMPFILE) && \
        rm $(COVER_TMPFILE).tmp; \
    done;

coverage: coverdata
	$(GO) tool cover -html=$(COVER_TMPFILE)
	@rm -f $(COVER_TMPFILE)

coveragefunc: coverdata
	$(GO) tool cover -func=$(COVER_TMPFILE)
	@rm -f $(COVER_TMPFILE)

# see 'go doc cmd/vet'
# 'go tool vet .' recursively descends the directory,
# vetting each package it finds.
vet:
	$(GOFOLDERS) | xargs $(GO) vet

# go get https://github.com/golang/lint
lint:
	lint=`$(GOFOLDERS) | xargs -L 1 $(GOLINT)`; if test -n "$$lint"; then echo "$$lint"; exit 1; fi

fmt:
	 $(GOFOLDERS) | xargs $(GO) fmt

doc:
	godoc -http=:6060 -index

clean:
	rm -f $(COVER_TMPFILE)
	$(GO) clean ./...

bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
		go get $$tool; \
	done

.PHONY: test unittest systest updatedeps coverdata coverage coveragefunc vet lint fmt doc clean bootstrap
