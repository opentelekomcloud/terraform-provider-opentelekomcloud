export PATH:=/usr/local/go/bin:$(PATH)
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
	@go test -i $(TEST) || exit 1
	@go test ./opentelekomcloud/acceptance/...

testacc: fmtcheck
	@TF_ACC=1 go test $(TEST) -v -timeout 720m

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


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile

