ifeq ($(OS), Windows_NT)
    GOOS_FLAG=windows
	BINARY_NAME=wireguard-tunnel.exe

    ifeq ($(PROCESSOR_ARCHITEW6432), AMD64)
        GOARCH_FLAG=amd64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        	GOARCH_FLAG=amd64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        	GOARCH_FLAG=x86
        endif
    endif
else
    UNAME_S := $(shell uname -s)
	BINARY_NAME=wireguard-tunnel

    ifeq ($(UNAME_S), Linux)
		GOOS_FLAG=linux
    endif
    ifeq ($(UNAME_S), Darwin)
		GOOS_FLAG=darwin
    endif
    
	UNAME_P := $(shell uname -p)
    
	ifeq ($(UNAME_P),x86_64)
        GOARCH_FLAG=amd64
    endif
    ifneq ($(filter %86, $(UNAME_P)),)
        GOARCH_FLAG=x86
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        GOARCH_FLAG=arm64
    endif
endif

all: bin/wireguard-tunnel-linux-amd64 bin/wireguard-tunnel-linux-arm64

bin/wireguard-tunnel-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o bin/wireguard-tunnel-linux-amd64 main.go

bin/wireguard-tunnel-linux-arm64:
	GOARCH=arm64 GOOS=linux go build -o bin/wireguard-tunnel-linux-arm64 main.go

build: 
	GOARCH=${GOARCH_FLAG} GOOS=${GOOS_FLAG} go build -o bin/${BINARY_NAME} main.go

clean:
	go clean
	rm -rf bin/

.PHONY: all build clean
