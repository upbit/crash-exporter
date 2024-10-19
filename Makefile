TARGET := crash_exporter
LDFLAGS := 

# REV := $(git rev-parse HEAD)
# CHANGES := $(test -n "$$(git status --porcelain)" && echo '+CHANGES' || true)

all: imports fmt lint build

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binaries, artifacts and releases.'
	@echo '    clean-artifacts    Remove build artifacts only.'
	@echo '    clean-releases     Remove releases only.'
	@echo '    clean-vendor       Remove content of the vendor directory.'
	@echo '    tools              Install tools needed by the project.'
	@echo '    test               Run unit tests.'
	@echo '    coverage           Report code tests coverage.'
	@echo '    vet                Run go vet.'
	@echo '    errors             Run errcheck.'
	@echo '    lint               Run golint.'
	@echo '    imports            Run goimports.'
	@echo '    fmt                Run go fmt.'
	@echo '    env                Display Go environment.'
	@echo '    build              Build project for current platform.'
	@echo '    build-all          Build project for all supported platforms.'
	@echo '    doc                Start Go documentation server on port 8080.'
	@echo '    release            Package and sing project for release.'
	@echo '    package-release    Package release and compress artifacts.'
	@echo '    sign-release       Sign release and generate checksums.'
	@echo '    check              Verify compiled binary.'
	@echo '    vendor             Update and save project build time dependencies.'
	@echo '    version            Display Go version.'
	@echo ''
	@echo 'Targets run by default are: imports, fmt, lint, vet, errors and build.'
	@echo ''

lint:
	@golangci-lint run ./...

vet:
	go vet -v ./...

imports:
	goimports -l -w .

fmt:
	go fmt ./...

clean:
	go clean

build:
	go build -v -ldflags "$(LDFLAGS)" -o "$(TARGET)" .

test:
	go test -v ./...

coverage: deps
	gocov test ./... > $(CURDIR)/coverage.out 2>/dev/null
	gocov report $(CURDIR)/coverage.out
	if test -z "$$CI"; then \
	  gocov-html $(CURDIR)/coverage.out > $(CURDIR)/coverage.html; \
	  if which open &>/dev/null; then \
	    open $(CURDIR)/coverage.html; \
	  fi; \
	fi

doc:
	godoc -http=:8080 -index

tools:
	# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
	go get golang.org/x/tools/cmd/goimports
	go get github.com/axw/gocov/gocov
	go get github.com/matm/gocov-html
