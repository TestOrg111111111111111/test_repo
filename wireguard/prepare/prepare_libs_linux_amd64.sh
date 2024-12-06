rm -rf bin

cd wireguard-go/
make
cd ../

cd wireguard-tools/src/
make
cd ../../

cd amneziawg-go/
make
cd ../

cd amneziawg-tools/src/
make
cd ../../

cd wireguard-windows/
build.bat
cd ../

mkdir bin/libs
cp wireguard-go/wireguard-go bin/libs/
cp wireguard-tools/src/wg bin/libs/wg
cp -r wireguard-tools/src/wg-quick bin/libs/wg-quick
cp amneziawg-go/amneziawg-go bin/libs/
cp amneziawg-tools/src/wg bin/libs/awg
cp -r amneziawg-tools/src/wg-quick bin/libs/awg-quick
cp -r wireguard-windows/amd64/ bin/libs/amd64/
cp -r wireguard-windows/arm64/ bin/libs/arm64/
cp -r wireguard-windows/x86/ bin/libs/x86/
