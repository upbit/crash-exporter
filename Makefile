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
	@echo '    tools              Install tools needed by the project.'
	@echo '    test               Run unit tests.'
	@echo '    coverage           Report code tests coverage.'
	@echo '    lint               Run golint.'
	@echo '    imports            Run goimports.'
	@echo '    fmt                Run go fmt.'
	@echo '    build              Build project for current platform.'
	@echo '    doc                Start Go documentation server on port 8080.'
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
	rm -f $(TARGET).*
	@rm -f coverage.out coverage.html

build:
	GOOS=linux GOARCH=amd64 go build -v -ldflags "$(LDFLAGS)" -o "$(TARGET).linux-amd64" .
	GOOS=linux GOARCH=arm64 go build -v -ldflags "$(LDFLAGS)" -o "$(TARGET).linux-arm64" .

package:
	mv $(TARGET).linux-amd64 $(TARGET) && tar czfv $(TARGET).linux-amd64.tar.gz $(TARGET) && rm $(TARGET)
	mv $(TARGET).linux-arm64 $(TARGET) && tar czfv $(TARGET).linux-arm64.tar.gz $(TARGET) && rm $(TARGET)

test:
	go test -v ./...

coverage:
	gocov test ./... > $(CURDIR)/coverage.out 2>/dev/null
	gocov report $(CURDIR)/coverage.out
	if test -z "$$CI"; then \
	  gocov-html $(CURDIR)/coverage.out > $(CURDIR)/coverage.html; \
	fi

doc:
	godoc -http=:8080 -index

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin latest
	go install golang.org/x/tools/cmd/godoc@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/matm/gocov-html/cmd/gocov-html@latest
