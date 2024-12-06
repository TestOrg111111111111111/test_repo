package main

import (
	"log"
	
	"os"
	"os/exec"
)

const WG_CONFIG_FOLDER = "C:\\Program Files\\DobbyVPN\\Config\\"
const WIREGUARD_EXE = "libs/amd64/wireguard.exe"

func tunnelWireguardOn(connectionData ConnectionData) {
	if err := os.MkdirAll(WG_CONFIG_FOLDER, os.ModePerm) ; err != nil {
		log.Printf("Error during %s folder create: %s\n", WG_CONFIG_FOLDER, err)
		os.Exit(1)
	}

	configPath := WG_CONFIG_FOLDER + connectionData.interfaceName + ".conf"
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)

	if err != nil {
		log.Printf("Error during %s file open: %s\n", configPath, err)
		os.Exit(1)
	}

	if _, err := file.WriteString(connectionData.configContent) ; err != nil {
		log.Printf("Error during %s file write: %s\n", configPath, err)
		os.Exit(1)
	}

	file.Close()

	cmd := exec.Command(WIREGUARD_EXE, "/installtunnelservice", configPath)

	if err := cmd.Run() ; err != nil {
		log.Printf("Error during wireguard tunnel set up: %s\n", err)
		os.Exit(1)
	} else {
		log.Println("installtunnelservice success")
	}
}

func tunnelWireguardOff(connectionData ConnectionData) {
	cmd := exec.Command(WIREGUARD_EXE, "/uninstalltunnelservice", connectionData.interfaceName)

	if err := cmd.Run() ; err != nil {
		log.Printf("Error during wireguard tunnel set down: %s\n", err)
		os.Exit(1)
	} else {
		log.Println("uninstalltunnelservice success")
	}
}

func tunnelAmneziaWGOn(connectionData ConnectionData) {
	log.Fatal("Unsupported")
}

func tunnelAmneziaWGOff(connectionData ConnectionData) {
	log.Fatal("Unsupported")
}
