#
GO := go
GOTEST_FLAGS :=
TEST_TIMEOUT := 120s

# Code Coverage Related
COV_FILE := coverage.txt
COV_FUNC_OUT := coverage_func.out
COVER_MODE := count

# External Tools
EXTERNAL_TOOLS=\
	golang.org/x/tools/cmd/cover

PKGS := $(shell go list ./... | grep -v /examples)

ifndef GOBIN
	GOBIN = $(GOPATH)/bin
endif

GOLINT := $(GOBIN)/golint

all: install

install:
	$(GO) install $(PKGS)

test: unittest vet

systest:
	$(GO) test $(PKGS) $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) -run SystemTest$

unittest:
	$(GO) test $(PKGS) $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) -run UnitTest$

updatedeps:
	$(GO) get -u github.com/mitchellh/mapstructure
	$(GO) get -u github.com/vaughan0/go-ini

coverdata:
	@$(GO) tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		$(GO) get -u golang.org/x/tools/cmd/cover; \
	fi
	@echo 'mode: $(COVER_MODE)' > $(COV_FILE)
	@for dir in $(PKGS); do \
		$(GO) test -covermode=$(COVER_MODE) -coverprofile=cov_tmp.out -run UnitTest$  $$dir || exit; \
		tail -n +2 cov_tmp.out >> $(COV_FILE) && \
		rm -f cov_tmp.out; \
	done;

coverage: coverdata
	$(GO) tool cover -html=$(COV_FILE)
	@rm -f $(COV_FILE)

coveragefunc: coverdata
	$(GO) tool cover -func=$(COV_FILE)
	#$(GO) tool cover -func=$(COV_FILE) > $(COV_FUNC_OUT)
	@rm -f $(COV_FILE)

# see 'go doc cmd/vet'
# 'go tool vet .' recursively descends the directory,
# vetting each package it finds.
vet:
	@$(GO) tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		$(GO) get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet ."
	@$(GO) tool vet . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Please check the reported files and constructs for errors "; \
		echo "and fix before submitting the code for review."; \
	fi

# go get https://github.com/golang/lint
lint:
	$(GOLINT) ./...

fmt:
	$(GO) fmt ./...

doc:
	godoc -http=:6060 -index

clean:
	rm -f $(COV_FILE) $(COV_FUNC_OUT)

bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
		go get $$tool; \
	done

.PHONY: test unittest systest updatedeps coverdata coverage coveragefunc vet lint fmt doc clean bootstrap
