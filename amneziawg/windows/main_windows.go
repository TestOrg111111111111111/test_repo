package main

import (
	"debug/pe"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/windows"

	"github.com/amnezia-vpn/amneziawg-windows/tunnel"

	"github.com/amnezia-vpn/amneziawg-windows-client/elevate"
	"github.com/amnezia-vpn/amneziawg-windows-client/l18n"
	"github.com/amnezia-vpn/amneziawg-windows-client/manager"
	"github.com/amnezia-vpn/amneziawg-windows-client/ringlogger"
)

func setLogFile() {
	logHandle, err := windows.GetStdHandle(windows.STD_ERROR_HANDLE)
	if logHandle == 0 || err != nil {
		logHandle, err = windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	}
	if logHandle == 0 || err != nil {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.NewFile(uintptr(logHandle), "stderr"))
	}
}

func fatal(v ...any) {
	log.Fatal(append([]any{l18n.Sprintf("Error: ")}, v...))
}

func fatalf(format string, v ...any) {
	fatal(l18n.Sprintf(format, v...))
}

func checkForWow64() {
	b, err := func() (bool, error) {
		var processMachine, nativeMachine uint16
		err := windows.IsWow64Process2(windows.CurrentProcess(), &processMachine, &nativeMachine)
		if err == nil {
			return processMachine != pe.IMAGE_FILE_MACHINE_UNKNOWN, nil
		}
		if !errors.Is(err, windows.ERROR_PROC_NOT_FOUND) {
			return false, err
		}
		var b bool
		err = windows.IsWow64Process(windows.CurrentProcess(), &b)
		if err != nil {
			return false, err
		}
		return b, nil
	}()
	if err != nil {
		fatalf("Unable to determine whether the process is running under WOW64: %v", err)
	}
	if b {
		fatalf("You must use the native version of AmneziaWG on this computer.")
	}
}

func checkForAdminGroup() {
	// This is not a security check, but rather a user-confusion one.
	var processToken windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY|windows.TOKEN_DUPLICATE, &processToken)
	if err != nil {
		fatalf("Unable to open current process token: %v", err)
	}
	defer processToken.Close()
	if !elevate.TokenIsElevatedOrElevatable(processToken) {
		fatalf("AmneziaWG may only be used by users who are a member of the Builtin %s group.", elevate.AdminGroupName())
	}
}

func checkForAdminDesktop() {
	adminDesktop, err := elevate.IsAdminDesktop()
	if !adminDesktop && err == nil {
		fatalf("AmneziaWG is running, but the UI is only accessible from desktops of the Builtin %s group.", elevate.AdminGroupName())
	}
}

func pipeFromHandleArgument(handleStr string) (*os.File, error) {
	handleInt, err := strconv.ParseUint(handleStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(handleInt), "pipe"), nil
}

func printUsage() {
	flags := [...]string{
		"/installtunnelservice CONFIG_PATH",
		"/uninstalltunnelservice TUNNEL_NAME",
		"/dumplog",
	}

	builder := strings.Builder{}
	for _, flag := range flags {
		builder.WriteString(fmt.Sprintf("    %s\n", flag))
	}
	fmt.Printf("Usage: %s [\n%s]\n", os.Args[0], builder.String())
	os.Exit(1)
}

func main() {
	if windows.SetDllDirectory("") != nil || windows.SetDefaultDllDirectories(windows.LOAD_LIBRARY_SEARCH_SYSTEM32) != nil {
		panic("failed to restrict dll search path")
	}

	if len(os.Args) < 2 {
		printUsage()
	}

	setLogFile()
	checkForWow64()

	switch os.Args[1] {
	case "/installtunnelservice":
		if len(os.Args) != 3 {
			printUsage()
		}

		err := manager.InstallTunnel(os.Args[2])
		if err != nil {
			fatal(err)
		}
		return
	case "/uninstalltunnelservice":
		if len(os.Args) != 3 {
			printUsage()
		}

		err := manager.UninstallTunnel(os.Args[2])
		if err != nil {
			fatal(err)
		}
		return
	case "/dumplog":
		if len(os.Args) != 2 {
			printUsage()
		}

		outputHandle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
		if err != nil {
			fatal(err)
		}
		if outputHandle == 0 {
			fatal("Stdout must be set")
		}
		file := os.NewFile(uintptr(outputHandle), "stdout")
		defer file.Close()
		logPath, err := manager.LogFile(false)
		if err != nil {
			fatal(err)
		}
		err = ringlogger.DumpTo(logPath, file, false)
		if err != nil {
			fatal(err)
		}
		return
	// For the inner service usage
	case "/tunnelservice":
		err := tunnel.Run(os.Args[2])
		if err != nil {
			fatal(err)
		}
		return
	}
}
