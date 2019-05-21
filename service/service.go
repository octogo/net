package service

import (
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/octogo/dict"
	"github.com/octogo/log"
	"github.com/octogo/rate"
)

// ConnectionHandler is a function that can handle a net.Conn and returns an
// error indicating why the connection terminated.
type ConnectionHandler func(net.Conn) error

// Service implements a service manager for a net.Listener.
type Service struct {
	// Listener holds the underlying net.Listener
	Listener net.Listener

	// Handler holds the designated handler function for new connections.
	Handler ConnectionHandler

	// MaxConnections defines the maximum number of connected clients.
	// Once this maximum has been reached, new clients will receive an
	// errTooManyConnections. A value of <=0 means unlimited connections.
	MaxConnections int64

	// MaxConnectionsPerIP defines the maximum number of connections per source
	// IP address. Once this maximum has been reached, new clients for the same
	// IP will receive an errTooManyConnectionsFromIP. A value of <=0 means
	// unlimited connections.
	MaxConnectionsPerIP int64

	// ConnectRate defines the maximum rate at which incoming connections
	// will be handled by this service. A value of nil disables rate-limiting.
	ConnectRate time.Duration

	// ConnectCount defines the maximum number of connections that are allowed
	// within ConnectRate.
	ConnectCount int

	// ConnectBurst defines the burst limit for incoming connections. A value
	// of <=0 disables the burst bucket.
	ConnectBurst int

	gauge           *rate.Gauge
	connections     *dict.Dict
	connectionCount int64
}

// New returns an initialized Service.
func New(l net.Listener, h ConnectionHandler) *Service {
	return &Service{
		Listener: l,
		Handler:  h,

		connections: dict.New(nil),
	}
}

// Run blocks and runs the service.
func (s *Service) Run() {
	s.gauge = rate.NewGauge(s.ConnectBurst, s.ConnectCount, s.ConnectRate)
	for {
		// obey the rate-limiting
		s.gauge.Wait()

		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println("listener terminated:", err)
			return
		}
		if conn == nil {
			continue
		}

		// handle the connection
		s.handleConnection(conn)
	}
}

func (s *Service) handleConnection(conn net.Conn) {
	defer conn.Close()

	// add client
	err := s.addClient(conn)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	// pass to ConnectionHandler func
	go func(c net.Conn, s *Service) {
		defer func(c net.Conn, s *Service) {
			// remove client
			err = s.delClient(c)
			if err != nil {
				log.Println(err)
			}
		}(conn, s)

		out := make(chan error)
		go func() {
			out <- s.Handler(conn)
		}()
		err := <-out
		if err != nil {
			log.Println("client failed:", err)
		}
	}(conn, s)
}

func (s *Service) addClient(conn net.Conn) error {
	// check max connections
	if s.MaxConnections > 0 {
		if s.connectionCount >= s.MaxConnections {
			return errTooManyConnections
		}
	}

	// obtain net.Addr from conn
	addr, err := net.ResolveTCPAddr(conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	// get or create a new []net.Conn for this connection's remote IP.
	ip := addr.IP.String()
	conns, found := s.connections.Get(ip)
	if !found {
		// create a new []net.Conn with a capacity of s.Config.MaxConnections and save it
		conns = dict.New(nil)
		s.connections.Set(ip, conns)
	}

	curConns, _ := conns.(*dict.Dict)

	// check max connections per IP
	if s.MaxConnectionsPerIP > 0 {
		if int64(len(curConns.Map())) >= s.MaxConnectionsPerIP {
			return errTooManyConnectionsPerIP
		}
	}

	curConns.Set(conn.RemoteAddr().String(), conn)
	atomic.AddInt64(&s.connectionCount, 1)
	return nil
}

func (s *Service) delClient(conn net.Conn) error {
	// obtain net.Addr from conn
	addr, err := net.ResolveTCPAddr(conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	// get []net.Conn for the connection's remote IP
	ip := addr.IP.String()
	conns, found := s.connections.Get(ip)
	if !found {
		return errors.New("no connections from this IP: " + ip)
	}

	conns.(*dict.Dict).Del(conn.RemoteAddr().String())
	if len(conns.(*dict.Dict).Map()) == 0 {
		s.connections.Del(ip)
	}

	s.connectionCount--
	return nil
}

// ConnectionCount returns the current connection count.
func (s Service) ConnectionCount() int64 {
	return s.connectionCount
}

// Measure returns the current rate of the gauge's speed-o-meter.
func (s Service) Measure() float64 {
	return s.gauge.Measure()
}
