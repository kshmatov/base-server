package handler

import (
	"github.com/kshmatov/base-server/socket"
	"github.com/kshmatov/base-server/entity"
	"github.com/kshmatov/base-server/logger"
	"sync"
)

type PacketHandler struct {
	sock socket.Socket
	entity entity.Entity
	logger logger.Logger
}

func Init(s socket.Socket, e entity.Entity, l logger.Logger) (*PacketHandler, error){
	ph := PacketHandler{
		sock: s,
		entity: e,
		logger: l,
	}
	return &ph, nil
}

func listenSocket(p PacketHandler, buffSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	d := make([]byte, buffSize)
	for {
		l, err := p.sock.Read(d)
		if err != nil {
			p.logger("ERROR", "On ReadSocket: %v", err)
			return
		}
		if l <= 0{
			p.logger("ERROR", "Readed %v bytes from socket", err)
			return
		}

	}
}