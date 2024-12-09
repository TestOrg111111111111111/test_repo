# AmneziaWG client

A simple script, that runs [AmneziaWG](https://docs.amnezia.org/documentation/amnezia-wg/) with the provided configuration

## Build

### Build on the current platform

```bash
make build
```

### Build specific platform:

```bash
make bin/amneziawg-tunnel-linux-amd64
```

## Usage

Usage is platform dependent

### Linux

```bash
amneziawg-tunnel-___ installtunnelservice CONFIG_PATH INTERFACE_NAME
```

This command runs AmneziaWG tunnel with the configuration provided in the `CONFIG_PATH` file.
This tunnel will be linked to the `INTERFACE_NAME` interface.
Tunnel log will be printed to the **STDOUT**

While this script is running, it will be implement AmneziaWG tunnel and 
[IPC configuration](https://www.wireguard.com/xplatform/) so that it can be additionally configured using `awg` and other utilities.

### Windows

#### Install tunnel

```bash
amneziawg-tunnel-___.exe /installtunnelservice CONFIG_ABSOLUTE_PATH
```

Runs tunnel service on the separate process. `CONFIG_ABSOLUTE_PATH` file must be with `${INTERFACE_NAME}.conf` format.

#### Remove tunnel

```bash
amneziawg-tunnel-___.exe /uninstalltunnelservice INTERFACE_NAME
```

Removed tunnel with the provided interface name

#### Get log

```bash
amneziawg-tunnel-___.exe /dumplog
```

Prints tunnel log to the **STDOUT**

## Config file format

```
[Interface]
PrivateKey = <...>
Address = <...>
MTU = <...>
DNS = <...>
Jc = <...>
Jmin = <...>
Jmax = <...>
S1 = <...>
S2 = <...>
H1 = <...>
H2 = <...>
H3 = <...>
H4 = <...>

[Peer]
PublicKey = <...>
AllowedIPs = <...>, <...>, <...>, <...>, ...
Endpoint = <...>
PersistentKeepalive = <...>

[Peer]
PublicKey = <...>
AllowedIPs = <...>, <...>, <...>, <...>, ...
Endpoint = <...>
PersistentKeepalive = <...>

[Peer]
...
```

## Additional documentation

* [Configure AmneziaWG on Unix server](./SERVER_CONFIG_UNIX.md)
