ifeq ($(OS), Windows_NT)
	GOOS_FLAG=windows

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

	BUILD_TARGET=bin/amneziawg-tunnel-${GOOS_FLAG}-${GOARCH_FLAG}.exe
else
	UNAME_S := $(shell uname -s)

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

	BUILD_TARGET=bin/amneziawg-tunnel-${GOOS_FLAG}-${GOARCH_FLAG}
endif

all: build

bin/amneziawg-tunnel-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o bin/amneziawg-tunnel-linux-amd64 main_linux.go

bin/amneziawg-tunnel-linux-arm64:
	GOARCH=arm64 GOOS=linux go build -o bin/amneziawg-tunnel-linux-arm64 main_linux.go

bin/amneziawg-tunnel-windows-amd64.exe:
	mkdir bin/
	cp wintun/bin/amd64/wintun.dll bin/
	GOARCH=amd64 GOOS=windows go build -o bin/amneziawg-tunnel-windows-amd64.exe main_windows.go

bin/amneziawg-tunnel-windows-arm64.exe:
	mkdir bin/
	cp wintun/bin/arm64/wintun.dll bin/
	GOARCH=arm64 GOOS=windows go build -o bin/amneziawg-tunnel-windows-arm64.exe main_windows.go

bin/amneziawg-tunnel-windows-x86.exe:
	mkdir bin/
	cp wintun/bin/x86/wintun.dll bin/
	GOARCH=x86 GOOS=windows go build -o bin/amneziawg-tunnel-windows-x86.exe main_windows.go

build: ${BUILD_TARGET}

clean:
	go clean
	rm -rf bin/

.PHONY: all build clean
