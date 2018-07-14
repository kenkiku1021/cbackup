GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=cbackup
BINARY_UNIX=$(BINARY_NAME)_unix
CGO_ENABLED=1

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
build-win:
	GOOS=windows $(GOBUILD)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
