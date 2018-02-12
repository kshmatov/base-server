package main

import (
	"net"
	"github.com/kshmatov/base-server/logger"
	"time"
	"fmt"
)

func listener(c <-chan bool) {
	l := logger.GetBaseConsoleLogger("SERVER")

	s,err := net.Listen("tcp", "localhost:9988")
	if err != nil {
		l("ERROR", "On Listen: %v", err)
		return
	}

	fin := false

	go func(){
		<- c
		fin = true
		s.Close()
		l("INFO", "Close listener")
	}()

	for {
		c, err := s.Accept()
		if err != nil {
			if fin {
				return
			} else {
				l("ERROR", "on Accept: %v", err)
				continue
			}
		}
		go func(c *net.TCPConn) {
			defer func(){
				l("TRACE", "RemoteAddr: %v (%T)", c.RemoteAddr(), c.RemoteAddr())
				f, err := c.File()
				l("TRACE", "FD: %v (%T), err: %v", f, f, err)
			}()

			defer c.Close()

			b := make([]byte, 1024)
			c.SetReadDeadline(time.Now().Add(time.Second * 5))
			i, err := c.Read(b)
			if err != nil {
				l("ERROR", "on Read: %v, %v", i, err)
				return
			}
			_, err = c.Write([]byte(fmt.Sprintf("Ok: %v bytes", i)))
			if err != nil {
				l("ERROR", "on Write: %v", err)
			}
			return
		}(c.(*net.TCPConn))

	}
}

func main() {
	ch := make(chan bool)
	defer close(ch)

	go listener(ch)

	time.Sleep(time.Second)

	l := logger.GetBaseConsoleLogger("DIALER")
	s, err := net.Dial("tcp", "localhost:9988")
	if err != nil {
		l("ERROR", "on 	Dial: %v", err)
		return
	}
	b := make([]byte, 1024)

	s.Write([]byte("Hello!"))

	i, err := s.Read(b)
	l("INFO", string(b[0:i]))
	s.Close()
	l("TRACE", "RemoteAddr: %v (%T)", s.RemoteAddr(), s.RemoteAddr())
	f, err :=s.(*net.TCPConn).File()
	l("TRACE", "FD: %v (%T), err: %v", f, f, err)

	time.Sleep(time.Minute * 5)
}
