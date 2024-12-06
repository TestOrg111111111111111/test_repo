package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const AWG = "./libs/awg"
const AWG_QUICK = "./libs/awg-quick/linux.bash"
const AWG_GO = "./libs/amneziawg-go"
const AWG_CONFIG_FOLDER = "/etc/amnezia/amneziawg/"
const WG = "./libs/wg"
const WG_QUICK = "./libs/wg-quick/linux.bash"
const WG_GO = "./libs/wireguard-go"
const WG_CONFIG_FOLDER = "/etc/wireguard/"

func runWireguardInterface(interfaceName string) error {
	var cmd = exec.Command("sudo", WG_GO, interfaceName)
	return cmd.Run()
}

func runAmneziaWGInterface(interfaceName string) error {
	var cmd = exec.Command("sudo", AWG_GO, interfaceName)
	return cmd.Run()
}

func setConfig(configContent string, interfaceName string, configFolder string, configScript string) error {
	configPath := configFolder + interfaceName + ".conf"
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	if _, err := file.WriteString(configContent); err != nil {
		return err
	}

	file.Close()

	cmd := exec.Command("sudo", configScript, "setconf", interfaceName, configPath)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func setWireguardConfig(configContent string, interfaceName string) error {
	return setConfig(configContent, interfaceName, WG_CONFIG_FOLDER, WG)
}

func setAmneziaWGConfig(configContent string, interfaceName string) error {
	return setConfig(configContent, interfaceName, AWG_CONFIG_FOLDER, AWG)
}

func addAddress(interfaceName string, address string) error {
	var cmd = exec.Command("sudo", "ip", "-4", "address", "add", address, "dev", interfaceName)
	return cmd.Run()
}

func setMtu(interfaceName string, mtu string) error {
	var cmd = exec.Command("sudo", "ip", "link", "set", "mtu", mtu, "up", "dev", interfaceName)
	return cmd.Run()
}

func addWireguardRoute(interfaceName string, address string) error {
	var table = "51820"

	if err := exec.Command("sudo", WG, "set", interfaceName, "fwmark", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "rule", "add", "not", "fwmark", table, "table", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "rule", "add", "table", "main", "suppress_prefixlength", "0").Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "route", "add", address, "dev", interfaceName, "table", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "sysctl", "-q", "net.ipv4.conf.all.src_valid_mark=1").Run(); err != nil {
		return err
	}

	return nil
}

func addAmneziaWGRoute(interfaceName string, address string) error {
	var table = "51820"

	if err := exec.Command("sudo", AWG, "set", interfaceName, "fwmark", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "rule", "add", "not", "fwmark", table, "table", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "rule", "add", "table", "main", "suppress_prefixlength", "0").Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "ip", "-4", "route", "add", address, "dev", interfaceName, "table", table).Run(); err != nil {
		return err
	}

	if err := exec.Command("sudo", "sysctl", "-q", "net.ipv4.conf.all.src_valid_mark=1").Run(); err != nil {
		return err
	}

	return nil
}

func tunnelWireguardOn(connectionData ConnectionData) {
	if err := runWireguardInterface(connectionData.interfaceName); err != nil {
		log.Fatalf("Wireguard interface run failed: %s\n", err)
	} else {
		log.Println("Wireguard interface run succeed")
	}

	if err := setWireguardConfig(connectionData.configContent, connectionData.interfaceName); err != nil {
		tunnelWireguardOff(connectionData.interfaceName)
		log.Fatalf("Wireguard interface config failed: %s\n", err)
	} else {
		log.Println("Wireguard interface config succeed")
	}

	if err := addAddress(connectionData.interfaceName, connectionData.address); err != nil {
		tunnelWireguardOff(connectionData.interfaceName)
		log.Fatalf("Wireguard interface address addition failed: %s\n", err)
	} else {
		log.Println("Wireguard interface address addition succeed")
	}

	if err := setMtu(connectionData.interfaceName, connectionData.mtu); err != nil {
		tunnelWireguardOff(connectionData.interfaceName)
		log.Fatalf("Wireguard interface mtu set failed: %s\n", err)
	} else {
		log.Println("Wireguard interface mtu set succeed")
	}

	// Set dns

	for _, aip := range strings.Split(connectionData.allowed_ip, " ") {
		if err := addWireguardRoute(connectionData.interfaceName, aip); err != nil {
			tunnelWireguardOff(connectionData.interfaceName)
			log.Fatalf("Wireguard interface route %s failed: %s\n", aip, err)
		} else {
			log.Printf("Wireguard interface route %s succeed\n", aip)
		}
	}
}

func tunnelWireguardOff(interfaceName string) {
	var cmd = exec.Command("sudo", "ip", "link", "del", interfaceName)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal("Could not remove Wireguard interface:\n error: %s\n output: %s\n", err, output)
	} else {
		log.Println("Removed Wireguard interface")
	}
}

func tunnelAmneziaWGOn(connectionData ConnectionData) {
	if err := runAmneziaWGInterface(connectionData.interfaceName); err != nil {
		log.Fatalf("AmneziaWG interface run failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface run succeed")
	}

	if err := setAmneziaWGConfig(connectionData.configContent, connectionData.interfaceName); err != nil {
		tunnelAmneziaWGOff(connectionData.interfaceName)
		log.Fatalf("AmneziaWG interface config failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface config succeed")
	}

	if err := addAddress(connectionData.interfaceName, connectionData.address); err != nil {
		tunnelAmneziaWGOff(connectionData.interfaceName)
		log.Fatalf("AmneziaWG interface address addition failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface address addition succeed")
	}

	if err := setMtu(connectionData.interfaceName, connectionData.mtu); err != nil {
		tunnelAmneziaWGOff(connectionData.interfaceName)
		log.Fatalf("AmneziaWG interface mtu set failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface mtu set succeed")
	}

	// Set dns

	for _, aip := range strings.Split(connectionData.allowed_ip, " ") {
		if err := addAmneziaWGRoute(connectionData.interfaceName, aip); err != nil {
			tunnelWireguardOff(connectionData.interfaceName)
			log.Fatalf("AmneziaWG interface route %s failed: %s\n", aip, err)
		} else {
			log.Printf("AmneziaWG interface route %s succeed\n", aip)
		}
	}
}

func tunnelAmneziaWGOff(interfaceName string) {
	var cmd = exec.Command("sudo", "ip", "link", "del", interfaceName)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatal("Could not remove AmneziaWG interface:\n error: %s\n output: %s\n", err, output)
	} else {
		log.Println("Removed AmneziaWG interface")
	}
}
