package main

import (
	"fmt"
	"os"
	"os/exec"
)

const WG_CONFIG_FOLDER = "C:\\Program Files\\DobbyVPN\\Config\\"

func (wg Wireguard) turn_on(interface_name string) {
	if err := os.MkdirAll(WG_CONFIG_FOLDER, os.ModePerm) ; err != nil {
		fmt.Printf("Error during %s folder create: %s\n", WG_CONFIG_FOLDER, err)
		os.Exit(1)
	}

	config_path := WG_CONFIG_FOLDER + interface_name + ".conf"
	file, err := os.OpenFile(config_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)

	if err != nil {
		fmt.Printf("Error during %s file open: %s\n", config_path, err)
		os.Exit(1)
	}

	if _, err := file.WriteString(wg.wg_content) ; err != nil {
		fmt.Printf("Error during %s file write: %s\n", config_path, err)
		os.Exit(1)
	}

	file.Close()

	cmd := exec.Command("libs/amd64/wireguard.exe", "/installtunnelservice", config_path)

	if err := cmd.Run() ; err != nil {
		fmt.Printf("Error during wireguard tunnel set up: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("installtunnelservice success")
	}
}

func (wg Wireguard) turn_off(interface_name string) {
	cmd := exec.Command("libs/amd64/wireguard.exe", "/uninstalltunnelservice", interface_name)

	if err := cmd.Run() ; err != nil {
		fmt.Printf("Error during wireguard tunnel set down: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("uninstalltunnelservice success")
	}
}
