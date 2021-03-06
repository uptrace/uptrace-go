ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

test:
	go test ./...
	go test ./... -short -race
	go test ./... -run=NONE -bench=. -benchmem
	env GOOS=linux GOARCH=386 go test ./...
	go vet ./...

tag:
	git tag $(VERSION)
	git tag extra/otellogrus/$(VERSION)
	git tag extra/otelzap/$(VERSION)

fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/uptrace/uptrace-go ./

go_mod_tidy:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy); \
	done
