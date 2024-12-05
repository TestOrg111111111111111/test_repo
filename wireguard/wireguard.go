package main
 
import (
    "fmt"
	"os"
	"os/exec"
)

type Wireguard struct {
	config_path string
}

func (wg Wireguard) run_interface(interface_name string) {
	var cmd = exec.Command("sudo", "./libs/wireguard-go", interface_name)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Could not run Wireguard interface: %s\n", err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Ran Wireguard interface\n")
	}
}

func (wg Wireguard) set_config(interface_name string) {
	var cmd = exec.Command("sudo", "./libs/wg", "setconf", interface_name, config_path)
	if output, err := cmd.CombinedOutput(); err != nil {		
		fmt.Printf("Could not configure Wireguard interface: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Set config\n")
	}
}

func (wg Wireguard) add_address(interface_name string, address string) {
	var cmd = exec.Command("sudo", "ip", "-4", "address", "add", address, "dev", interface_name)
	if err := cmd.Run(); err != nil {		
		fmt.Printf("Could not add address %s: %s\n", address, err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Add %s address\n", address)
	}
}

func (wg Wireguard) set_mtu(interface_name string, mtu string) {
	var cmd = exec.Command("sudo", "ip", "link", "set", "mtu", mtu, "up", "dev", interface_name)
	if err := cmd.Run(); err != nil {		
		fmt.Printf("Could not set up mtu %s: %s\n", mtu, err)
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Set %s mtu\n", mtu)
	}
}

func (wg Wireguard) turn_on(interface_name string) {
	wg.run_interface(interface_name)
	wg.set_config(interface_name)
	wg.add_address(interface_name, ADDRESS)
	wg.set_mtu(interface_name, MTU)
	// Set dns
	wg.add_address(interface_name, ALLOWED_IP)
}

func (wg Wireguard) turn_off(interface_name string) {
	var cmd = exec.Command("sudo", "ip", "link", "del", interface_name)
	if output, err := cmd.CombinedOutput(); err != nil {		
		fmt.Printf("Could not remove Wireguard interface: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		wg.turn_off(interface_name)
		os.Exit(1)
	} else {
		fmt.Printf("+ Removed interface\n")
	}
}