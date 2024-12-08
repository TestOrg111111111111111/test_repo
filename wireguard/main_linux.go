package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/amnezia-vpn/amneziawg-go/conn"
	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/ipc"
	"github.com/amnezia-vpn/amneziawg-go/tun"
)

func configureAmneziaWG(device *device.Device, configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}

	return device.IpcSetOperation(file)
}

func addAddress(interfaceName string, address string) error {
	var cmd = exec.Command("sudo", "ip", "-4", "address", "add", address, "dev", interfaceName)
	return cmd.Run()
}

func setMtu(interfaceName string, mtu string) error {
	var cmd = exec.Command("sudo", "ip", "link", "set", "mtu", mtu, "up", "dev", interfaceName)
	return cmd.Run()
}

func addAmneziaWGRoute(interfaceName string, address string) error {
	var table = "51820"

	// FIXME
	if err := exec.Command("sudo", "awg", "set", interfaceName, "fwmark", table).Run(); err != nil {
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
		log.Fatal("Could not remove AmneziaWG interface:\n error: %s\n output: %s\n", err, output)
	} else {
		log.Println("Removed AmneziaWG interface")
	}
}

func installTunnel(configPath string, interfaceName string) {
	// Logger definition:
	logLevel := device.LogLevelVerbose

	logger := device.NewLogger(
		logLevel,
		fmt.Sprintf("(%s) ", interfaceName),
	)

	// open TUN device (or use supplied fd)
	tdev, err := tun.CreateTUN(interfaceName, device.DefaultMTU)

	logger.Verbosef("Starting amneziawg")

	if err == nil {
		realInterfaceName, err2 := tdev.Name()
		if err2 == nil {
			interfaceName = realInterfaceName
		}
	} else {
		logger.Errorf("Failed to create TUN device: %v", err)
		os.Exit(1)
	}

	// open UAPI file
	fileUAPI, err := ipc.UAPIOpen(interfaceName)

	if err != nil {
		logger.Errorf("UAPI listen error: %v", err)
		os.Exit(1)
	}

	// Start device:

	device := device.NewDevice(tdev, conn.NewDefaultBind(), logger)

	logger.Verbosef("Device started")

	errs := make(chan error)
	term := make(chan os.Signal, 1)

	uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)

	if err != nil {
		logger.Errorf("Failed to listen on uapi socket: %v", err)
		os.Exit(1)
	}

	go func() {
		// Extra ipc configuration:
		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go device.IpcHandle(conn)
		}
	}()

	logger.Verbosef("UAPI listener started")

	// Configure AmneziaWG

	configureAmneziaWG(device, configPath)
	postConfigAmneziaWg(interfaceName)

	// wait for program to terminate

	signal.Notify(term, unix.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	case <-device.Wait():
	}

	// clean up

	tunnelAmneziaWGOff(interfaceName)
	uapi.Close()
	device.Close()

	logger.Verbosef("Shutting down")
}

func printUnixUsage() {
	fmt.Printf("Usage: %s [intalltunnelservice CONFIG_PATH INTERFACE_NAME]\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUnixUsage()
	}

	switch os.Args[1] {
	case "installtunnelservice":
		if len(os.Args) != 4 {
			printUnixUsage()
		}

		installTunnel(os.Args[2], os.Args[3])
		os.Exit(0)
	}

	printUnixUsage()
}
