package goConnectPool

import "net"

// Pool interface
type Pool interface {
	// Get in blocking mode
	Get() (net.Conn, error)

	// TryGet in nonBlocking mode
	TryGet() (net.Conn, error)

	// Close the Pool
	Close()

	// Len of the Pool
	Len() int

	// LenActives of the Pool
	LenActives() int
}
