package goConnectPool

import "net"

type Pool interface {
	Get() (net.Conn, error)

	TryGet() (net.Conn, error)

	Close()

	Len() int

	LenActives() int
}
