package server

import "fmt"

// Error while server is listening
type ServerListenError struct {
	host     string
	port     int
	protocol string
}

func (e *ServerListenError) Error() string {
	return fmt.Sprintf("Error while listening on %v:%d protocol %v", e.host, e.port, e.protocol)
}
