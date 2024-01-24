package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tcpserver/config"
	"time"

	"github.com/sirupsen/logrus"
)

type Server interface {
	Start() error
	IsRunning() bool
}

type server struct {
	info        ServerInfo
	running     bool
	connections []*net.TCPConn

	log *logrus.Entry
}

func New(serverInfo ServerInfo) *server {
	s := new(server)
	s.info = serverInfo
	s.running = false
	s.log = config.GetLogger("server").WithField("local", s.String())

	return s
}

func (s *server) Start() error {
	s.log.Info("Starting server")

	listener, err := s.getListener()
	if err != nil {
		s.log.WithField("cause", err).Error("Cannot start server")
		return err
	}
	defer listener.Close()
	s.running = true

	s.log.Info("Listening on ", listener.Addr().String())
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			s.log.WithField("cause", err).Error("error accepting connection")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *server) getListener() (*net.TCPListener, error) {
	ips, err := net.LookupIP(s.info.Host)
	if err != nil {
		s.log.Error("Error looking up host ", s.info.Host)
		return nil, err
	}

	var listener *net.TCPListener

	for _, ip := range ips {
		address := &net.TCPAddr{IP: ip, Port: s.info.Port}
		listener, err = net.ListenTCP(s.info.Protocol, address)
		if err != nil {
			s.log.Warning("Cannot listen to address ", address, ", cause: ", err)
			continue
		}
		break
	}
	if listener == nil {
		return nil, &ServerListenError{s.info.Host, s.info.Port, s.info.Protocol}
	}

	return listener, nil
}

func (s *server) handleConnection(conn *net.TCPConn) {
	remoteAddr := conn.RemoteAddr().String()
	s.log = s.log.WithField("remote", remoteAddr)
	s.log.Info("Connected to ", remoteAddr)

	s.addConnection(conn)
	defer s.removeConnection(conn)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			s.log.WithField("cause", err).Error("error reading data")
			return
		}

		trimmed := strings.TrimSpace(string(data))
		s.log.WithField("message", trimmed).Debug("handled new message")

		conn.Write([]byte("OK\n"))
	}
}

func (s *server) addConnection(conn *net.TCPConn) {
	s.connections = append(s.connections, conn)
}

func (s *server) removeConnection(conn *net.TCPConn) {
	var founds []int
	for i, e := range s.connections {
		if e == conn {
			founds = append(founds, i)
		}
	}
	if len(founds) == 0 {
		panic(fmt.Sprintf("cannot remove connection %v", conn))
	}
	for _, found := range founds {
		copy(s.connections[found:], s.connections[found+1:])
		s.connections = s.connections[:len(s.connections)-1]
	}
}

func (s *server) ListConnections() {
	for i, conn := range s.connections {
		s.log.Infof("connection[%d] = %v", i, conn.RemoteAddr().String())
	}
}

func (s *server) CloseAll() {
	total := len(s.connections)
	coco := make([]*net.TCPConn, total)
	copy(coco[:], s.connections)
	for i, conn := range coco {
		s.log.Infof("closing connection %d/%d : %v", i+1, total, conn.RemoteAddr().String())
		err := conn.Close()
		if err != nil {
			s.log.Errorf("Error closing connection %d/%d : %v", i+1, total, err)
		} else {
			s.log.Infof("Successful closed connection %d/%d", i+1, total)
		}
	}
	time.Sleep(2 * time.Second)
}

func (s *server) Connections() []*net.TCPConn {
	return s.connections
}

func (s *server) IsRunning() bool {
	return s.running
}

func (s *server) String() string {
	return fmt.Sprintf("%v://%v:%d", s.info.Protocol, s.info.Host, s.info.Port)
}
