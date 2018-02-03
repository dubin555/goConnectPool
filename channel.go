package goConnectPool

import (
	"net"
	"sync"
)

// The channel pool will implements the Pool
type channelPool struct {
	mu sync.Mutex

	// conns channel
	conns chan net.Conn

	// actives channel
	actives chan struct{}

	// factory to generate net.Conn
	factory Factory

	// channel to sync the close signal
	chanClosed chan bool

	// sign symbol to show closed state
	closed bool
}

// Factory function to generate a net.Conn
type Factory func() (net.Conn, error)

// NewChannelPool method to create factory function to new a channel pool
func NewChannelPool(initialCap int, maxCap int, maxActive int, factory Factory) (Pool, error) {

	// check capacity setting
	if initialCap < 0 || maxCap <= 0 || maxActive <= 0 || initialCap > maxCap {
		return nil, NewErrInvalidCap("Invalid capacity setting")
	}

	c := &channelPool{
		conns:      make(chan net.Conn, maxCap),
		actives:    make(chan struct{}, maxActive),
		factory:    factory,
		chanClosed: make(chan bool, 10),
	}

	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			c.Close()
			return nil, NewErrFactoryInitial("Factory initialize connect pool error")
		}
		c.conns <- conn
	}

	return c, nil
}

func (c *channelPool) getConnsAndFactory() (chan net.Conn, Factory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conns, c.factory
}

func (c *channelPool) getActives() chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	actives := c.actives
	return actives
}

func (c *channelPool) setStateClose() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	// will not block since the chanClosed will be taken one
	c.chanClosed <- true
}

// Get a net.Conn from the channelPool, if there is an available connection in channel,
// use it, Or create a new one.
// wrap it to PoolConn use release the source when closed.
func (c *channelPool) get(isNonBlocking bool) (net.Conn, error) {

	if c.closed {
		return nil, NewErrPoolClosed("Pool closed")
	}

	if isNonBlocking {
		select {
		case c.actives <- struct{}{}:
		// actives channel still have room
		case <-c.chanClosed:
			c.setStateClose()
			return nil, NewErrPoolClosed("Pool closed")
		default:
			return nil, NewErrConnLimit("conn limit reacted")
		}
	} else {
		select {
		case c.actives <- struct{}{}:
			// block when pool has room
		case <-c.chanClosed:
			c.setStateClose()
			return nil, NewErrPoolClosed("Pool closed")
		}
	}

	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, NewErrPoolClosed("Conns channel is null")
	}

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, NewErrPoolClosed("Conn is null")
		}
		return c.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		return c.wrapConn(conn), nil
	}
}

func (c *channelPool) Get() (net.Conn, error) {
	return c.get(false)
}

func (c *channelPool) TryGet() (net.Conn, error) {
	return c.get(true)
}

func (c *channelPool) put(conn net.Conn) error {
	if conn == nil {
		return NewErrPoolClosed("connection is null")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		return conn.Close()
	}

	// put the conn back into pool. If already full, just close the conn.
	select {
	case c.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	conns := c.conns
	if !c.closed {
		c.chanClosed <- true
	}
	// change point inside channelPool to nil, for GC
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		// already close
		return
	}

	// close the conns for rejecting the request
	close(conns)

	// close the conn in the channel
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) Len() int {
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}

func (c *channelPool) LenActives() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.actives)
}
