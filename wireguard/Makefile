BINARY_NAME=wireguard-tunnel

clean:
	go clean
	rm -rf libs/
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows.exe

build_unix:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go amneziawg-unix.go wireguard-unix.go

build_win:
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe main.go amneziawg-windows.go wireguard-windows-amd64.go

prepare_unix:
	sh ./prepare_libs_unix.sh

prepare_win:
	sh ./prepare_libs_windows.sh
