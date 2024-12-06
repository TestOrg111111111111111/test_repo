package main

import (
	"flag"
	"os"
	"log"
)

type ConnectionData struct{
	configContent   string
	interfaceName   string
	address         string
	dns             string
	mtu             string
	allowed_ip      string
}

const CONGIG_PATH_DEFAULT = "awg0.conf"
const CONGIG_PATH_USAGE = "Path to the config"
const TYPE_DEFAULT = "amneziawg"
const TYPE_USAGE = "Connection type, values: wireguard (alternative: wg) or amneziawg (alternative: awg))"
const INTERFACE_NAME_DEFAULT = "awg0"
const INTERFACE_NAME_USAGE = "Interface name"
const ADDRESS_DEFAULT = "10.9.9.2/24"
const ADDRESS_USAGE = "Interface address"
const DNS_DEFAULT = "8.8.8.8"
const DNS_USAGE = "Interface DNS"
const MTU_DEFAULT = "1420"
const MTU_USAGE = "Tunnel MTU"
const ALLOWED_IP_DEFAULT = "0.0.0.0/0"
const ALLOWED_IP_USAGE = "Interface peer allowed ips"
const MODE_DEFAULT = "yes"
const MODE_USAGE = "Turn on/Turn of (yes/no)"

func main() {
	var config_path string
	var connection_type string
	var interfaceName string
	var address string
	var dns string
	var mtu string
	var allowed_ip string
	var mode string

	flag.StringVar(&config_path, "config", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE)
	flag.StringVar(&connection_type, "type", TYPE_DEFAULT, TYPE_USAGE)
	flag.StringVar(&interfaceName, "iname", INTERFACE_NAME_DEFAULT, INTERFACE_NAME_USAGE)
	flag.StringVar(&address, "address", ADDRESS_DEFAULT, ADDRESS_USAGE)
	flag.StringVar(&dns, "dns", DNS_DEFAULT, DNS_USAGE)
	flag.StringVar(&mtu, "mtu", MTU_DEFAULT, MTU_USAGE)
	flag.StringVar(&allowed_ip, "ips", ALLOWED_IP_DEFAULT, ALLOWED_IP_USAGE)
	flag.StringVar(&mode, "mode", MODE_DEFAULT, MODE_USAGE)

	flag.Parse()
	
	data, err := os.ReadFile(config_path)
	if err != nil { log.Fatal(err) }
	configContent := string(data)
	connectionData := ConnectionData { configContent, interfaceName, address, dns, mtu, allowed_ip }

	switch mode {
	case "yes":
		switch connection_type {
		case "amneziawg", "awg":
			tunnelAmneziaWGOn(connectionData)
		case "wireguard", "wg":
			tunnelWireguardOn(connectionData)
		}
	case "no":
		switch connection_type {
		case "amneziawg", "awg":
			tunnelAmneziaWGOff(connectionData)
		case "wireguard", "wg":
			tunnelWireguardOff(connectionData)
		}
	}
}
