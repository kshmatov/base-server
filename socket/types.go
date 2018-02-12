package socket

import (
	"io"
	"github.com/kshmatov/base-server/logger"
)

type Socket interface {
	io.Reader
	io.Writer
	io.Closer
	SetId(uint32)
	SetLogger(logger logger.Logger)
}

type OnBuffer func(sender Socket, b []byte)(int, error)
type OnEvent func(sender Socket) error
