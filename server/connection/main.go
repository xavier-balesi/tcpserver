package connection

import "net"

type connection struct {
	net.Conn
}
