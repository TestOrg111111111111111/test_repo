rm -rf libs/

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

mkdir libs
cp wireguard-go/wireguard-go libs/
cp wireguard-tools/src/wg libs/wg
cp -r wireguard-tools/src/wg-quick libs/wg-quick
cp amneziawg-go/amneziawg-go libs/
cp amneziawg-tools/src/wg libs/awg
cp -r amneziawg-tools/src/wg-quick libs/awg-quick
