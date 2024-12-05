package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

type Connection interface {
	turn_on(interface_name string)
	turn_off(interface_name string)
}

type WireguardUnix struct {
	address    string
	dns        string
	mtu        string
	allowed_ip string
	wg_content string
}

type AmneziaWGUnix struct {
	address     string
	dns         string
	mtu         string
	allowed_ip  string
	awg_content string
}

const CONGIG_PATH_DEFAULT = "awg0.conf"
const CONGIG_PATH_USAGE = "Path to the config"
const TYPE_DEFAULT = "amneziawg"
const TYPE_USAGE = "Connection type, values: wireguard (alternative: wg) or amneziawg (alternative: awg))"
const ADDRESS_DEFAULT = "10.9.9.2/24"
const ADDRESS_USAGE = "Interface address"
const DNS_DEFAULT = "8.8.8.8"
const DNS_USAGE = "Interface DNS"
const MTU_DEFAULT = "1420"
const MTU_USAGE = "Tunnel MTU"
const ALLOWED_IP_DEFAULT = "0.0.0.0/0"
const ALLOWED_IP_USAGE = "Interface peer allowed ips"

var config_path string
var connection_type string
var address string
var dns string
var mtu string
var allowed_ip string

func test_connection(conn Connection, name string) {
	conn.turn_on(name)
	conn.turn_off(name)
}

func build_wireguard_unix() WireguardUnix {
	data, err := os.ReadFile(config_path)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wg_content := string(data)

	return WireguardUnix{address, dns, mtu, allowed_ip, wg_content}
}

func build_wireguard() Connection {
	switch runtime.GOOS {
	case "windows":
		os.Exit(1)
	case "darwin":
		os.Exit(1)
	case "linux":
		return build_wireguard_unix()
	default:
		os.Exit(1)
	}

	return build_wireguard_unix()
}

func build_amneziawg_unix() AmneziaWGUnix {
	data, err := os.ReadFile(config_path)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wg_content := string(data)

	return AmneziaWGUnix{address, dns, mtu, allowed_ip, wg_content}
}

func build_amneziawg() Connection {
	switch runtime.GOOS {
	case "windows":
		os.Exit(1)
	case "darwin":
		os.Exit(1)
	case "linux":
		return build_amneziawg_unix()
	default:
		os.Exit(1)
	}

	return build_amneziawg_unix()
}

func main() {
	flag.StringVar(&config_path, "config", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE)
	flag.StringVar(&connection_type, "type", TYPE_DEFAULT, TYPE_USAGE)
	flag.StringVar(&address, "address", ADDRESS_DEFAULT, ADDRESS_USAGE)
	flag.StringVar(&dns, "dns", DNS_DEFAULT, DNS_USAGE)
	flag.StringVar(&mtu, "mtu", MTU_DEFAULT, MTU_USAGE)
	flag.StringVar(&allowed_ip, "ips", ALLOWED_IP_DEFAULT, ALLOWED_IP_USAGE)
	flag.Parse()

	var wg = build_wireguard()
	test_connection(wg, "wg0")

	var awg = build_amneziawg()
	test_connection(awg, "awg0")
}
