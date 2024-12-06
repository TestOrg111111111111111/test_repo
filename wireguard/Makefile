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

all: build

build: prepare
	GOARCH=${GOARCH_FLAG} GOOS=${GOOS_FLAG} go build -o bin/${BINARY_NAME} main.go amneziawg-${GOOS_FLAG}-${GOARCH_FLAG}.go wireguard-${GOOS_FLAG}-${GOARCH_FLAG}.go

prepare:
	sh ./prepare/prepare_libs_${GOOS_FLAG}_${GOARCH_FLAG}.sh

clean:
	go clean
	rm -rf bin/

.PHONY: prepare
