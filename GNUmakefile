SHELL=/bin/bash

export PATH:=/usr/local/go/bin:~/go/bin/:$(PATH)
TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')
PKG_NAME=opentelekomcloud

default: build

build: fmtcheck
	go install

release:
	goreleaser release

snapshot:
	goreleaser release --snapshot --parallelism 2 --rm-dist

test: fmtcheck
	go test -v ./...

acceptance: fmtcheck
	@bash "$(CURDIR)/scripts/run-acceptance.sh"

vet:
	@echo "go vet ."
	@go vet $$(go list ./...); if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

lint:
	golangci-lint run ./...

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST)

tools:
	@echo "==> installing required tooling..."
	go install github.com/katbyte/terrafmt@latest

tflint: tools
	./scripts/run-tflint.sh

tffmtfix: tools
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@find ./opentelekomcloud/acceptance -type f -name "*_test.go" | sort | while read f; do terrafmt fmt -f $$f; done
	@echo "==> Fixing docs terraform blocks code with terrafmt..."
	@find ./docs -type f -name "*.md" | sort | while read f; do terrafmt fmt $$f; done
	@echo "==> Fixing example docs terraform blocks code with terrafmt..."
	@find ./examples -type f -name "*.md" | sort | while read f; do terrafmt fmt $$f; done

.PHONY: build test acceptance vet fmt fmtcheck errcheck test-compile tflint tffmtfix lint
