package main
 
import (
    "fmt"
	"flag"
)

const ADDRESS = "10.9.9.2/24"
const MTU = "1420"
const ALLOWED_IP = "0.0.0.0/0"

type Connection interface {
	turn_on(interface_name string)
	turn_off(interface_name string)
}

const CONGIG_PATH_DEFAULT = "/etc/amnezia/amneziawg/"
const CONGIG_PATH_USAGE = "Path to the config"
const TYPE_DEFAULT = "amneziawg"
const TYPE_USAGE = "Connection type, values: wireguard (alternative: wg) or amneziawg (alternative: awg))"

var config_path string
var connection_type string

func test_connection(conn Connection, name string) {
	conn.turn_on(name)
	conn.turn_off(name)
}

func main() {
	flag.StringVar(&config_path, "c", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE + " (shorthand)")
	flag.StringVar(&config_path, "config", CONGIG_PATH_DEFAULT, CONGIG_PATH_USAGE)
	flag.StringVar(&connection_type, "t", TYPE_DEFAULT, TYPE_USAGE + " (shorthand)")
	flag.StringVar(&connection_type, "type", TYPE_DEFAULT, TYPE_USAGE)
	flag.Parse()
	
    fmt.Printf("Parameters: %s, %s\n", config_path, connection_type)

	var wg = Wireguard { config_path }
	test_connection(wg, "wg0")

	var awg = AmneziaWG { config_path }
	test_connection(awg, "awg0")
}
