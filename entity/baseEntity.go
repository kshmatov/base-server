package entity

import "github.com/kshmatov/base-server/packet"

type Entity interface{
	Handle(packet.BaseMessage) (error)
}