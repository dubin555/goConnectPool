package goConnectPool

// Error for invalid capacity of channel Pool.
type ErrInvalidCap struct {
	message string
}

func NewErrInvalidCap(message string) *ErrInvalidCap {
	return &ErrInvalidCap{
		message: message,
	}
}

func (e *ErrInvalidCap) Error() string {
	return e.message
}

// Error for factory initial the pool
type ErrFactoryInitial struct {
	message string
}

func NewErrFactoryInitial(message string) *ErrFactoryInitial {
	return &ErrFactoryInitial{
		message: message,
	}
}

func (e *ErrFactoryInitial) Error() string {
	return e.message
}

// Error for pool closed
type ErrPoolClosed struct {
	message string
}

func NewErrPoolClosed(message string) *ErrPoolClosed {
	return &ErrPoolClosed{
		message: message,
	}
}

func (e *ErrPoolClosed) Error() string {
	return e.message
}

// Error for connection limit
type ErrConnLimit struct {
	message string
}

func NewErrConnLimit(message string) *ErrConnLimit {
	return &ErrConnLimit{
		message: message,
	}
}

func (e *ErrConnLimit) Error() string {
	return e.message
}
