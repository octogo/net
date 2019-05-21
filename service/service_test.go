package service

import (
	"log"
	"net"
	"testing"
	"time"
)

func handler(c net.Conn) error {
	return nil
}

func TestConstructor(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Constructor was incorrect, got: %v", err)
		}
	}()
	l, err := net.Listen("tcp", "")
	if err != nil {
		t.Error(err)
	}
	s := New(l, handler)
	defer l.Close()

	s.ConnectBurst = 10
	s.ConnectCount = 10
	s.ConnectRate = time.Second
}

func ExampleNew() {
	handler := func(c net.Conn) error {
		_, err := c.Write([]byte(c.RemoteAddr().String()))
		return err
	}

	l, err := net.Listen("tcp", "")
	if err != nil {
		log.Fatal(err)
	}
	service := New(l, handler)
	service.Run()
}
