GO := go
GOTEST_FLAGS := -v
COV_FILE := coverage.out
COV_FUNC_OUT := coverage_func.out
GOPKG := github.com/aristanetworks/goeapi

ifndef GOBIN
	GOBIN = $(GOPATH)/bin
endif

GOLINT := $(GOBIN)/golint

test:
	$(GO) test $(GOTEST_FLAGS) ./...
	@$(MAKE) vet

testsys:
	$(GO) test ./ ./module $(GOTEST_FLAGS) -run SystemTest$

testunit:
	$(GO) test ./ ./module $(GOTEST_FLAGS) -run UnitTest$

updatedeps:
	$(GO) get -u github.com/mitchellh/mapstructure
	$(GO) get -u github.com/vaughan0/go-ini

coverage:
	@$(GO) tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		$(GO) get -u golang.org/x/tools/cmd/cover; \
	fi
	$(GO) test $(GOPKG) -coverprofile=$(COV_FILE)
	$(GO) test $(GOPKG)/module -coverprofile=coverage_api.out
	@tail -n +2 coverage_api.out >> $(COV_FILE)
	$(GO) tool cover -html=$(COV_FILE)
	$(GO) tool cover -func=$(COV_FILE) > $(COV_FUNC_OUT)
	@rm coverage_api.out

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

.PHONY: test testunit testsys updatedeps coverage vet lint fmt doc clean
