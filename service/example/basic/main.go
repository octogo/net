package main

import (
	"math/rand"
	"net"
	"time"

	"github.com/octogo/log"
	"github.com/octogo/net/service"
)

var logger = log.NewLogger("example")

func connectionHandler(conn net.Conn) error {
	logger := logger.NewLogger(conn.RemoteAddr().String())

	addr, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		logger.Fatal(err)
	}

	conn.Write([]byte(addr.String() + "\n"))
	<-time.After(time.Duration(rand.Intn(5) * int(time.Second)))
	return nil
}

func monitor(s *service.Service) {
	logger := logger.NewLogger("monitor")
	for {
		<-time.After(time.Second)
		logger.Noticef(
			"Monitor: %v %d %v",
			s.Measure(),
			s.ConnectionCount(),
		)
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:31337")
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Listening:", listener.Addr())

	s := service.New(listener, connectionHandler)
	s.ConnectBurst = 10
	s.ConnectCount = 10
	s.ConnectRate = time.Second

	go monitor(s)
	s.Run()
}
