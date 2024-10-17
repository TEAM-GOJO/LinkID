GOOS=linux
GOARCH=amd64

all: build

build:
	@echo "Compiling src/blockchain.go into ./blockchain"
	GOARCH=$(ARCH) GOOS=$(OS) go build -o blockchain src/blockchain.go

run:
	@echo "Running src/blockchain.go"
	go run src/blockchain.go

clean:
	@echo "Removing current records from /records"
	rm records/*