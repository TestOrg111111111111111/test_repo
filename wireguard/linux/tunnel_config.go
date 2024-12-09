package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/vishvananda/netlink"
)

func configureAmneziaWG(device *device.Device, configPath string, interfaceName string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	configString := string(configData)
	config, err := FromWgQuickWithUnknownEncoding(configString, interfaceName)
	if err != nil {
		return err
	}

	uapiString, err := config.ToUAPI()
	if err != nil {
		return err
	}

	log.Printf("Appliying uapi string:\n%s\n", uapiString)

	reader := strings.NewReader(uapiString)

	return device.IpcSetOperation(reader)
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

func setMtu(interfaceName string, mtu string) error {
	var cmd = exec.Command("sudo", "ip", "link", "set", "mtu", mtu, "up", "dev", interfaceName)
	return cmd.Run()
}

func addAmneziaWGRoute(interfaceName string, address string) error {
	var tableString = "51820"

	// sudo ip rule add not fwmark <table> table <table>
	if err := exec.Command("sudo", "ip", "rule", "add", "not", "fwmark", tableString, "table", tableString).Run(); err != nil {
		fmt.Println("#2")
		return err
	}

	// // sudo ip rule add table main suppress_prefixlength 0
	if err := exec.Command("sudo", "ip", "rule", "add", "table", "main", "suppress_prefixlength", "0").Run(); err != nil {
		fmt.Println("#3")
		return err
	}

	// sudo ip route add address dev <interfaceName> table <table>
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return err
	}

	_, dst, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}

	route := netlink.Route{LinkIndex: link.Attrs().Index, Dst: dst, Table: 51820}

	if err := netlink.RouteAdd(&route); err != nil {
		fmt.Println("#4")
		return err
	}

	// sudo sysctl -q net.ipv4.conf.all.src_valid_mark=1
	if err := exec.Command("sudo", "sysctl", "-q", "net.ipv4.conf.all.src_valid_mark=1").Run(); err != nil {
		fmt.Println("#5")
		return err
	}

	return nil
}

func postConfigAmneziaWg(interfaceName string) {
	if err := addAddress(interfaceName, "10.9.9.2/32"); err != nil {
		log.Fatalf("AmneziaWG interface address addition failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface address addition succeed")
	}

	if err := setMtu(interfaceName, "1420"); err != nil {
		log.Fatalf("AmneziaWG interface mtu set failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface mtu set succeed")
	}

	// Set dns

	if err := addAmneziaWGRoute(interfaceName, "0.0.0.0/0"); err != nil {
		log.Fatalf("AmneziaWG interface route %s failed: %s\n", "0.0.0.0/0", err)
	} else {
		log.Printf("AmneziaWG interface route %s succeed\n", "0.0.0.0/0")
	}
}

func tunnelAmneziaWGOff(interfaceName string) {
	var cmd = exec.Command("sudo", "ip", "link", "del", interfaceName)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Could not remove AmneziaWG interface:\n error: %s\n output: %s\n", err, output)
	} else {
		log.Println("Removed AmneziaWG interface")
	}
}
