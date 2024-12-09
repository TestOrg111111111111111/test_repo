package main

import (
	"fmt"
	"os"
)

func printUsage() {
	commands := `
	intalltunnelservice CONFIG_PATH INTERFACE_NAME
	install CONFIG_PATH INTERFACE_NAME
`
	fmt.Printf("Usage: %s [%s]\n", os.Args[0], commands)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	switch os.Args[1] {
	case "installtunnelservice", "install":
		if len(os.Args) != 4 {
			printUsage()
		}

		installTunnel(os.Args[2], os.Args[3])
		os.Exit(0)
	}

	printUsage()
}
