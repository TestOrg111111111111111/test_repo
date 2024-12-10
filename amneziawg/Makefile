ifeq ($(OS), Windows_NT)
	ifeq ($(PROCESSOR_ARCHITEW6432), AMD64)
		GOARCH_FLAG=amd64
		BUILD_TARGET=bin/windows/amd64/amneziawg-tunnel.exe
	else
		ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
			GOARCH_FLAG=amd64
			BUILD_TARGET=bin/windows/amd64/amneziawg-tunnel.exe
		endif
		ifeq ($(PROCESSOR_ARCHITECTURE),x86)
			GOARCH_FLAG=386
			BUILD_TARGET=bin/windows/x86/amneziawg-tunnel.exe
		endif
	endif
else
	BUILD_TARGET=bin/linux/amneziawg-tunnel
endif

build: ${BUILD_TARGET}

bin/:
	mkdir bin/
bin/windows: bin/
	mkdir bin/windows/
bin/linux: bin/
	mkdir bin/linux/
bin/windows/amd64: bin/windows
	mkdir bin/windows/amd64
bin/windows/x86: bin/windows
	mkdir bin/windows/x86

bin/linux/amneziawg-tunnel: bin/linux
	GOOS=linux go build -o bin/linux/amneziawg-tunnel linux/main.go linux/tunnel.go linux/amneziawg.go linux/config.go linux/parser.go linux/writer.go

bin/windows/amd64/amneziawg-tunnel.exe: bin/windows/amd64
	cp windows/wintun/bin/amd64/wintun.dll bin/windows/amd64/
	GOARCH=amd64 GOOS=windows go build -o bin/windows/amd64/amneziawg-tunnel.exe windows/main_windows.go

bin/windows/x86/amneziawg-tunnel.exe: bin/windows/x86
	cp windows/wintun/bin/x86/wintun.dll bin/windows/x86/
	GOARCH=386 GOOS=windows go build -o bin/windows/x86/amneziawg-tunnel.exe windows/main_windows.go

all: bin/linux/amneziawg-tunnel bin/windows/amd64/amneziawg-tunnel.exe bin/windows/x86/amneziawg-tunnel.exe

clean:
	go clean
	rm -rf bin/

.PHONY: build all clean
