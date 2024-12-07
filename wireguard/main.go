package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"os/exec"
	"os/signal"
	"golang.org/x/sys/unix"

	"github.com/amnezia-vpn/amneziawg-go/conn"
	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/ipc"
	"github.com/amnezia-vpn/amneziawg-go/tun"
)

const CONGIG_PATH_DEFAULT = "awg0.conf"
const CONGIG_PATH_USAGE = "Path to the config"
const LOG_LEVEL_DEFAULT = "debug"
const LOG_LEVEL_USAGE = "Log level"
var configPathVar string
var logLevelVar string
var interfaceName string

func configureAmneziaWG(device *device.Device, configPath string) error {
	file, err := os.Open(configPath)
	if err != nil { return err }

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

func postConfigAmneziaWg() {
	if err := addAddress("awg0", "10.9.9.2/32"); err != nil {
		log.Fatalf("AmneziaWG interface address addition failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface address addition succeed")
	}

	if err := setMtu("awg0", "1420"); err != nil {
		log.Fatalf("AmneziaWG interface mtu set failed: %s\n", err)
	} else {
		log.Println("AmneziaWG interface mtu set succeed")
	}

	// Set dns

	if err := addAmneziaWGRoute("awg0", "0.0.0.0/0"); err != nil {
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

func main() {
	flag.StringVar(&configPathVar, "config", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE)
	flag.StringVar(&logLevelVar, "log", LOG_LEVEL_DEFAULT, LOG_LEVEL_USAGE)
	flag.StringVar(&interfaceName, "iname", "awg0", "...")
	flag.Parse()

	// get log level

	logLevel := func() int {
		switch logLevelVar {
		case "verbose", "debug":
			return device.LogLevelVerbose
		case "error":
			return device.LogLevelError
		case "silent":
			return device.LogLevelSilent
		}
		return device.LogLevelError
	}()

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

	configureAmneziaWG(device, configPathVar)
	postConfigAmneziaWg()

	// wait for program to terminate

	signal.Notify(term, unix.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	case <-device.Wait():
	}

	// clean up

	tunnelAmneziaWGOff("awg0")
	uapi.Close()
	device.Close()

	logger.Verbosef("Shutting down")
}
