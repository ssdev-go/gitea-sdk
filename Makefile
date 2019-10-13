.PHONY: all
all: clean test build

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: fmt
fmt:
	find . -name "*.go" -type f ! -path "./vendor/*" ! -path "./benchmark/*" | xargs gofmt -s -w

.PHONY: vet
vet:
	cd gitea && go vet ./...

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u golang.org/x/lint/golint; \
	fi
	cd gitea && golint -set_exit_status

.PHONY: test
test:
	cd gitea && go test -cover -coverprofile coverage.out

.PHONY: bench
bench:
	cd gitea && go test -run=XXXXXX -benchtime=10s -bench=. || exit 1

.PHONY: build
build:
	cd gitea && go build
