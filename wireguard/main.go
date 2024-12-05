package main
 
import (
    "fmt"
	"flag"
	"os"
	"os/exec"
)

const ADDRESS = "10.9.9.2/24"
const MTU = "1420"
const ALLOWED_IP = "0.0.0.0/0"

type connection interface {
	turn_on(interface_name string)
	turn_off(interface_name string)
}

type Wireguard struct {
	config_path string
}

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

// Wireguard:

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

const CONGIG_PATH_DEFAULT = "/etc/amnezia/amneziawg/"
const CONGIG_PATH_USAGE = "Path to the config"
const TYPE_DEFAULT = "amneziawg"
const TYPE_USAGE = "Connection type, values: wireguard (alternative: wg) or amneziawg (alternative: awg))"

var config_path string
var connection_type string

func main() {
	flag.StringVar(&config_path, "c", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE + " (shorthand)")
	flag.StringVar(&config_path, "config", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE)
	flag.StringVar(&connection_type, "t", TYPE_DEFAULT, TYPE_USAGE + " (shorthand)")
	flag.StringVar(&connection_type, "type", TYPE_DEFAULT, TYPE_USAGE)
	flag.Parse()
	
    fmt.Printf("Parameters: %s, %s\n", config_path, connection_type)

	var wireguard = Wireguard { config_path }
	wireguard.turn_on("wg0")
	wireguard.turn_off("wg0")
}
