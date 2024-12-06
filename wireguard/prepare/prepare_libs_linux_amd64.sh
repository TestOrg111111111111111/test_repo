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

mkdir bin
mkdir bin/libs
cp wireguard-go/wireguard-go bin/libs/
cp wireguard-tools/src/wg bin/libs/wg
cp -r wireguard-tools/src/wg-quick bin/libs/wg-quick
cp amneziawg-go/amneziawg-go bin/libs/
cp amneziawg-tools/src/wg bin/libs/awg
cp -r amneziawg-tools/src/wg-quick bin/libs/awg-quick
