GOOS=linux
GOARCH=amd64

all: build

build:
	@echo "Compiling src/main.go into ./main..."
	GOARCH=$(ARCH) GOOS=$(OS) go build -o main src/main.go
