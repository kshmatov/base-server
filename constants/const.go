package constants

import "time"

const (
	TCPKeepAlivePeriod  = time.Second * 30
	TCPBaseDeadline = time.Second * 20
	TCPAwaitNewData = time.Minute * 30
)