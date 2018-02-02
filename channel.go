package goConnectPool

import (
	"sync"
	"net"
)

type channelPool struct {
	mu sync.Mutex

	// conns channel
	conns chan net.Conn

	// actives channel
	actives chan struct{}

	// factory to generate net.Conn
	factory Factory

	chanClosed chan bool

	closed bool
}

type Factory func() (net.Conn, error)

func NewChannelPool(initialCap int, maxCap int, maxActive int, factory Factory) (Pool, error) {

	// check capacity setting
	if initialCap < 0 || maxCap <= 0 || maxActive <= 0 || initialCap > maxCap {
		return nil, NewErrInvalidCap("Invalid capacity setting")
	}

	c := &channelPool{
		conns: make(chan net.Conn, maxCap),
		actives:make(chan struct{}, maxActive),
		factory: factory,
		chanClosed:make(chan bool, 10),
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

func (c *channelPool) tryClose(conn net.Conn) error {

	// actives minus one when closing the connection
	// would not be blocked because actives plus one when open it.
	<- c.actives

	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (c *channelPool) getConnsAndFactory() (chan net.Conn, Factory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conns, c.factory
}

func (c * channelPool) getActives() chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	actives := c.actives
	return actives
}

// Get a net.Conn from the channelPool, if there is an available connection in channel,
// use it, Or create a new one.
// wrap it to PoolConn use release the source when closed.
func (c *channelPool) get(isNonBlocking bool) (net.Conn, error) {

	if isNonBlocking {
		select {
		case c.actives <- struct{}{}:
		// actives channel still have room
		case <- c.chanClosed:
			c.mu.Lock()
			c.closed = true
			c.mu.Unlock()
			return nil, NewErrPoolClosed("Pool closed")
		default:
			return nil, NewErrConnLimit("conn limit reacted")
		}
	} else {
		select {
		case c.actives <- struct {}{}:
			// block when pool has room
		case <- c.chanClosed:
			c.mu.Lock()
			c.closed = true
			c.mu.Unlock()
			return nil, NewErrPoolClosed("Pool closed")
		}
	}

	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, NewErrPoolClosed("Conns channel is null")
	}

	select {
	case conn := <- conns:
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
	actives := c.actives
	if !c.closed {
		c.chanClosed <- true
	}
	// change point inside channelPool to nil, for GC
	c.conns = nil
	c.actives = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		// already close
		return
	}

	// close the conns for rejecting the request
	close(conns)

	// close the actives
	close(actives)

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