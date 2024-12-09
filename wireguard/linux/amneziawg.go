package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	sysctl "github.com/lorenzosaino/go-sysctl"
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

	if config.setUpInterface(); err != nil {
		return err
	} else {
		log.Println("Interface set up success")
	}

	for _, address := range config.Interface.Addresses {
		if err := config.addAddress(address.String()); err != nil {
			return err
		} else {
			log.Printf("Address %s addition success\n", address.String())
		}
	}

	if err := config.setDns(); err != nil {
		return err
	}

	for _, peer := range config.Peers {
		for _, allowed_ip := range peer.AllowedIPs {
			if err := config.addRoute(allowed_ip.String()); err != nil {
				log.Printf("Route from %s address to %s is failed", allowed_ip.String(), config.Name)
				return err
			} else {
				log.Printf("Route from %s address to %s is successful", allowed_ip.String(), config.Name)
			}
		}
	}

	return err
}

func (config *Config) setUpInterface() error {
	link, err := netlink.LinkByName(config.Name)
	if err != nil {
		return err
	}

	return netlink.LinkSetUp(link)
}

func (config *Config) addAddress(address string) error {
	// sudo ip -4 address add <address> dev <interfaceName>
	link, err := netlink.LinkByName(config.Name)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	return netlink.AddrAdd(link, addr)
}

func (config *Config) setDns() error {
	// TODO
	return nil
}

func (config *Config) addRoute(address string) error {
	// sudo ip rule add not fwmark <table> table <table>
	ruleNot := netlink.NewRule()
	ruleNot.Invert = true
	ruleNot.Mark = uint32(config.Interface.FwMark)
	ruleNot.Table = int(config.Interface.FwMark)
	if err := netlink.RuleAdd(ruleNot); err != nil {
		return err
	}

	// sudo ip rule add table main suppress_prefixlength 0
	ruleAdd := netlink.NewRule()
	ruleAdd.Table = unix.RT_TABLE_MAIN
	ruleAdd.SuppressPrefixlen = 0
	if err := netlink.RuleAdd(ruleAdd); err != nil {
		return err
	}

	// sudo ip route add <address> dev <interfaceName> table <table>
	link, err := netlink.LinkByName(config.Name)
	if err != nil {
		return err
	}

	_, dst, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}

	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: dst, Table: 51820}

	if err := netlink.RouteAdd(&route); err != nil {
		return err
	}

	// sudo sysctl -q net.ipv4.conf.all.src_valid_mark=1
	if err := sysctl.Set("net.ipv4.conf.all.src_valid_mark", "1"); err != nil {
		return err
	}

	return nil
}

func (config *Config) tunnelOff() {
	link, err := netlink.LinkByName(config.Name)
	if err == nil {
		err = netlink.LinkDel(link)
	}

	if err != nil {
		log.Fatalf("Could not remove AmneziaWG interface: %s\n", err)
	} else {
		log.Println("Removed AmneziaWG interface")
	}
}
