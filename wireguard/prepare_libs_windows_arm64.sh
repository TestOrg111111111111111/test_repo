rm -rf libs/

wireguard-windows/build.bat

mkdir libs
cp -r wireguard-windows/amd64/ libs/amd64/
cp -r wireguard-windows/arm64/ libs/arm64/
cp -r wireguard-windows/x86/ libs/x86/
