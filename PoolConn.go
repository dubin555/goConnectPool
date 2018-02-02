package goConnectPool

import (
	"net"
	"sync"
)

// PoolConn is a wrapper of net.Conn to change the default Close() method.
type PoolConn struct {
	net.Conn
	mu sync.Mutex
	c *channelPool
	unusable bool
}

// Put the resource back into pool is possible.
func (p *PoolConn) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.unusable {
		if p.Conn != nil {
			// close the conn
			// actives minus one when close the connection
			<-p.c.actives
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.put(p.Conn)
}

func (p *PoolConn) setUnusable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.unusable = true
}

// Wrapper function
func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}