BINARY_NAME=wireguard-tunnel

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go amneziawg-unix.go wireguard-unix.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go amneziawg-windows.go wireguard-windows.go

test: build
	sudo ./${BINARY_NAME}-linux --ips="0.0.0.0/0 0.0.0.0/1 0.0.0.0/2"

clean:
	go clean
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows
