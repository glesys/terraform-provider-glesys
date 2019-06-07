GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: build

build: fmtcheck
	go install
vet:
	@echo "go vet ."
	@go vet $$(go list) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

test: fmt vet

.PHONY: build vet fmt fmtcheck test
