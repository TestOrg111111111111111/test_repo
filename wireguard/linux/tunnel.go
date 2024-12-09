package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/amnezia-vpn/amneziawg-go/conn"
	"github.com/amnezia-vpn/amneziawg-go/device"
	"github.com/amnezia-vpn/amneziawg-go/ipc"
	"github.com/amnezia-vpn/amneziawg-go/tun"
	"golang.org/x/sys/unix"
)

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
