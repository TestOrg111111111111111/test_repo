package main
 
import (
    "fmt"
	"os"
	"os/exec"
)

type AmneziaWG struct {
	config_path string
}

func (wg AmneziaWG) run_interface(interface_name string) {
	var cmd = exec.Command("sudo", "./libs/amneziawg-go", interface_name)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not run AmneziaWG interface: %s\n", err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Ran AmneziaWG interface\n")
	}
}

func (wg AmneziaWG) set_config(interface_name string) {
	var cmd = exec.Command("sudo", "./libs/awg", "setconf", interface_name, config_path)
	if output, err := cmd.CombinedOutput(); err != nil {		
		fmt.Printf("Could not configure AmneziaWG interface: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Set config\n")
	}
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
	wg.add_address(interface_name, ADDRESS)
	wg.set_mtu(interface_name, MTU)
	// Set dns
	wg.add_address(interface_name, ALLOWED_IP)
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
