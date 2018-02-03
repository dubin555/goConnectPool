package goConnectPool

// ErrInvalidCap for invalid capacity of channel Pool.
type ErrInvalidCap struct {
	message string
}

// NewErrInvalidCap to create a new ErrInvalidCap
func NewErrInvalidCap(message string) *ErrInvalidCap {
	return &ErrInvalidCap{
		message: message,
	}
}

func (e *ErrInvalidCap) Error() string {
	return e.message
}

// ErrFactoryInitial for factory initial the pool
type ErrFactoryInitial struct {
	message string
}

// NewErrFactoryInitial to create a new ErrFactoryInitial
func NewErrFactoryInitial(message string) *ErrFactoryInitial {
	return &ErrFactoryInitial{
		message: message,
	}
}

func (e *ErrFactoryInitial) Error() string {
	return e.message
}

// ErrPoolClosed for pool closed
type ErrPoolClosed struct {
	message string
}

// NewErrPoolClosed to create a new ErrPoolClosed
func NewErrPoolClosed(message string) *ErrPoolClosed {
	return &ErrPoolClosed{
		message: message,
	}
}

func (e *ErrPoolClosed) Error() string {
	return e.message
}

// ErrConnLimit for connection limit
type ErrConnLimit struct {
	message string
}

// NewErrConnLimit to create a new ErrConnLimit
func NewErrConnLimit(message string) *ErrConnLimit {
	return &ErrConnLimit{
		message: message,
	}
}

func (e *ErrConnLimit) Error() string {
	return e.message
}
