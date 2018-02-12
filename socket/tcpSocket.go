package socket

import (
	"net"
	"github.com/kshmatov/base-server/logger"
	"github.com/kshmatov/base-server/constants"
	"time"
	"github.com/pkg/errors"
)

type TCPSocket struct {
	c            *net.TCPConn
	uid          uint32
	log          logger.Logger
	readTimeout  time.Duration
	writeTimeout time.Duration
	onRead       OnBuffer
	onClose      OnEvent
	closed       bool
}

func NewTCPSocket(c *net.TCPConn, log logger.Logger) *TCPSocket {
	s := TCPSocket{
		c: c,
		log: log,
	}
	s.c.SetKeepAlive(true)
	s.c.SetKeepAlivePeriod(constants.TCPKeepAlivePeriod)
	return &s
}

func (t *TCPSocket)SetReadTimeout(to time.Duration) {
	t.readTimeout = to
}

func (t *TCPSocket)SetWriteTimeout(to time.Duration) {
	t.writeTimeout = to
}

func (t *TCPSocket)Read(buf []byte)(int, error) {
	t.c.SetReadDeadline(time.Now().Add(t.readTimeout))
	i, err := t.c.Read(buf)
	if err != nil || i <=0 {
		return i, err
	}
	if t.onRead != nil {
		return t.onRead(t, buf)
	}
	return i, err
}

func (t *TCPSocket) Write(buf []byte)(int, error) {
	t.c.SetWriteDeadline(time.Now().Add(t.writeTimeout))
	return t.c.Write(buf)
}

func (t *TCPSocket) SetId(id uint32) {
	t.uid = id
}

func (t *TCPSocket) SetLogger(l logger.Logger) {
	t.log = l
}

func (t *TCPSocket) Close() error {
	var err error
	if t.closed {
		return errors.New("Socket is closed already")
	}

	err = t.c.Close()
	f, _ := t.c.File()
	t.closed = f != nil

	if t.onClose != nil && t.closed{
		err = t.onClose(t)
	}

	return err
}

func (t *TCPSocket) IsClosed()bool {
	return t.closed
}

func (t *TCPSocket) GetId() uint32{
	return t.uid
}