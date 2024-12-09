package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/vishvananda/netlink"
)

func ConfigFromPath(configPath string, interfaceName string) (*Config, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	return FromWgQuickWithUnknownEncoding(string(configData), interfaceName)
}

func (config *Config) tunnelOn(device *device.Device) error {
	uapiString, err := config.ToUAPI()
	if err != nil {
		return err
	}

	log.Printf("Appliying uapi string:\n%s\n", uapiString)

	reader := strings.NewReader(uapiString)

	if err := device.IpcSetOperation(reader); err != nil {
		return nil
	} else {
		log.Println("IPC config set success")
	}

	if setUpInterface(config.Name); err != nil {
		return err
	} else {
		log.Println("Interface set up success")
	}

	for _, address := range config.Interface.Addresses {
		if err := addAddress(config.Name, address.String()); err != nil {
			return err
		} else {
			log.Printf("Address %s addition success\n", address.String())
		}
	}

	if err := setDns(config.Name, config.Interface.DNS); err != nil {
		return err
	}

	for _, peer := range config.Peers {
		for _, allowed_ip := range peer.AllowedIPs {
			if err := addRoute(config.Name, allowed_ip.String()); err != nil {
				log.Printf("Route from %s address to %s is failed", allowed_ip.String(), config.Name)
				return err
			} else {
				log.Printf("Route from %s address to %s is successful", allowed_ip.String(), config.Name)
			}
		}
	}

	return err
}

func setUpInterface(interfaceName string) error {
	var cmd = exec.Command("sudo", "ip", "link", "set", "up", "dev", interfaceName)
	return cmd.Run()
}

func addAddress(interfaceName string, address string) error {
	// sudo ip -4 address add <address> dev <interfaceName>
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	return netlink.AddrAdd(link, addr)
}

func setDns(interfaceName string, dns []net.IP) error {
	// TODO
	return nil
}

func addRoute(interfaceName string, address string) error {
	var tableString = "51820"

	// sudo ip rule add not fwmark <table> table <table>
	if err := exec.Command("sudo", "ip", "rule", "add", "not", "fwmark", tableString, "table", tableString).Run(); err != nil {
		return err
	}

	// sudo ip rule add table main suppress_prefixlength 0
	if err := exec.Command("sudo", "ip", "rule", "add", "table", "main", "suppress_prefixlength", "0").Run(); err != nil {
		return err
	}

	// sudo ip route add <address> dev <interfaceName> table <table>
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}

	_, dst, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return err
	}

	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: dst, Table: 51820}

	if err := netlink.RouteAdd(&route); err != nil {
		return err
	}

	// sudo sysctl -q net.ipv4.conf.all.src_valid_mark=1
	if err := exec.Command("sudo", "sysctl", "-q", "net.ipv4.conf.all.src_valid_mark=1").Run(); err != nil {
		return err
	}

	return nil
}

func (config *Config) tunnelOff() {
	var cmd = exec.Command("sudo", "ip", "link", "del", config.Name)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Could not remove AmneziaWG interface:\n error: %s\n output: %s\n", err, output)
	} else {
		log.Println("Removed AmneziaWG interface")
	}
}
