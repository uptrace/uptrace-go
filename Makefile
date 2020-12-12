ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

all:
	go test ./...
	go test ./... -short -race
	go test ./... -run=NONE -bench=. -benchmem
	env GOOS=linux GOARCH=386 go test ./...
	go vet ./...
	golangci-lint run

tag:
	git tag $(VERSION)
	git tag extra/otellogrus/$(VERSION)

go_mod_tidy:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go mod tidy); \
	done
