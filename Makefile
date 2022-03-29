PACKAGE_DIRS := $(shell find . -mindepth 1 -type f -name 'go.mod' -exec dirname {} \; | sort)

test:
	go test ./...
	go test ./... -short -race
	go test ./... -run=NONE -bench=. -benchmem
	env GOOS=linux GOARCH=386 go test ./...
	go vet ./...

fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/uptrace/uptrace-go ./

go_mod_tidy:
	set -e; for dir in $(PACKAGE_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy -compat=1.17); \
	done
