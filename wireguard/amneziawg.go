package main

import (
	"fmt"
	"os"
	"os/exec"
)

const AWG = "./libs/awg"
const AWG_QUICK = "./libs/awg-quick/linux.bash"
const AWG_GO = "./libs/amneziawg-go"
const AWG_CONFIG_FOLDER = "/etc/amnezia/amneziawg/"

func (wg AmneziaWG) run_interface(interface_name string) {
	var cmd = exec.Command("sudo", AWG_GO, interface_name)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not run AmneziaWG interface: %s\n", err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Ran AmneziaWG interface\n")
	}
}

func (wg AmneziaWG) set_config(interface_name string) {
	file_path := AWG_CONFIG_FOLDER + interface_name + ".conf"
	file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE, 0604)

	if err != nil {
		fmt.Printf("Could not configure AmneziaWG interface: %s\n", err)
		wg.turn_off(interface_name)
		os.Exit(1)
	}

	if _, err := file.WriteString(wg.awg_content); err != nil {
		fmt.Printf("Could not configure AmneziaWG interface: %s\n", err)
		wg.turn_off(interface_name)
		os.Exit(1)
	}

	file.Close()
	cmd := exec.Command("sudo", AWG, "setconf", interface_name, file_path)

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Could not configure AmneziaWG interface: %s\n", err)
		fmt.Printf("=========\n%s\n", string(output))
		wg.turn_off(interface_name)
		fmt.Printf("=========\n")
		os.Exit(1)
	}

	fmt.Printf("+ Set config\n")
}

func (wg AmneziaWG) add_address(interface_name string, address string) {
	var cmd = exec.Command("sudo", "ip", "-4", "address", "add", address, "dev", interface_name)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not add address %s: %s\n", address, err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Add %s address\n", address)
	}
}

func (wg AmneziaWG) set_mtu(interface_name string, mtu string) {
	var cmd = exec.Command("sudo", "ip", "link", "set", "mtu", mtu, "up", "dev", interface_name)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not set up mtu %s: %s\n", mtu, err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Set %s mtu\n", mtu)
	}
}

func (wg AmneziaWG) turn_on(interface_name string) {
	wg.run_interface(interface_name)
	wg.set_config(interface_name)
	wg.add_address(interface_name, wg.address)
	wg.set_mtu(interface_name, wg.mtu)
	// Set dns
	wg.add_address(interface_name, wg.allowed_ip)
}

func (wg AmneziaWG) turn_off(interface_name string) {
	var cmd = exec.Command("sudo", "ip", "link", "del", interface_name)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Could not remove AmneziaWG interface: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Removed interface\n")
	}
}
