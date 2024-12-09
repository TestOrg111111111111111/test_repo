# Configure AmneziaWG on Unix server

## Preparation

### Download required libraries:

```bash
git clone https://github.com/amnezia-vpn/amneziawg-go
git clone https://github.com/amnezia-vpn/amneziawg-tools
```

### Build libraries:

Build `amneziawg-go/amneziawg-go`:

```bash
cd amneziawg-go/
make
cd ../
```

Build `amneziawg-tools/src/wg-quick/awg`:

```bash
cd amneziawg-tools/src/
make
cd ../../
```

## Intall tunnel

Variables:
* `INTERFACE_NAME` - just a simple interface name
* `CONFIG_PATH` - path the the server config file, its format can be seen a the end of file
* `SERVER_ADDRESS` - Server tunnel addres, just use `10.9.9.1/24`

### Run interface

```bash
amneziawg-go/amneziawg-go INTERFACE_NAME
```

### Configure interface

```bash
amneziawg-tools/src/wg-quick/awg setconf INTERFACE_NAME CONFIG_PATH
ip -4 address add SERVER_ADDRESS dev INTERFACE_NAME
ip link set mtu 1420 up dev INTERFACE_NAME
```

### Routing config

```bash
sysctl -w net.ipv4.ip_forward=1
```

And then `PostUp` iptables command, that can be generated using [awgcfg.py](https://gist.githubusercontent.com/remittor/8c3d9ff293b2ba4b13c367cc1a69f9eb/raw/awgcfg.py) util or using the next pattern:
```
iptables -A INPUT -p udp --dport <SERVER_PORT> -m conntrack --ctstate NEW -j ACCEPT --wait 10 --wait-interval 50; iptables -A FORWARD -i <SERVER_IFACE> -o <SERVER_TUN> -j ACCEPT --wait 10 --wait-interval 50; iptables -A FORWARD -i <SERVER_TUN> -j ACCEPT --wait 10 --wait-interval 50; iptables -t nat -A POSTROUTING -o <SERVER_IFACE> -j MASQUERADE --wait 10 --wait-interval 50; ip6tables -A FORWARD -i <SERVER_TUN> -j ACCEPT --wait 10 --wait-interval 50; ip6tables -t nat -A POSTROUTING -o <SERVER_IFACE> -j MASQUERADE --wait 10 --wait-interval 50
```

#### Config generating using [awgcfg.py](https://gist.githubusercontent.com/remittor/8c3d9ff293b2ba4b13c367cc1a69f9eb/raw/awgcfg.py)

```bash
wget -O awgcfg.py https://gist.githubusercontent.com/remittor/8c3d9ff293b2ba4b13c367cc1a69f9eb/raw/awgcfg.py
python3 awgcfg.py --make CONFIG_PATH -i SERVER_ADDRESS -p SERVER_PORT
python3 awgcfg.py --create
python3 awgcfg.py -a "ClientName"
python3 awgcfg.py -c
```

This generates server config to the `CONFIG_PATH` file and client config to the `ClientName.conf` path.

> **Notice I**: There is a small bug at this script, it generates invalid obfuscation parameter (jc, jmin etc), so that copy that fields from one config to another

> **Notice II**: This script generates server config with extra parameters: `Address`, `PostUp` and `PostDown`. Theese parameters should be removed, PostUp command should be run after interface configuration (*Routing config* step), PostDown command should be run after interface shutdown and Address parameter should be used as SERVER_ADDRESS parameter 

## Conclusion:

```bash
git clone https://github.com/amnezia-vpn/amneziawg-go
git clone https://github.com/amnezia-vpn/amneziawg-tools

cd amneziawg-go/
make
cd ../

cd amneziawg-tools/src/
make
cd ../../

amneziawg-go/amneziawg-go INTERFACE_NAME

amneziawg-tools/src/wg-quick/awg setconf INTERFACE_NAME CONFIG_PATH
ip -4 address add SERVER_ADDRESS dev INTERFACE_NAME
ip link set mtu 1420 up dev INTERFACE_NAME

sysctl -w net.ipv4.ip_forward=1
./PostUp.sh
```
