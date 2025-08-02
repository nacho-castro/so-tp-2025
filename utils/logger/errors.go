package logger

import "errors"

var (
	ErrBadRequest   = errors.New("bad Request")
	ErrIsEmpty      = errors.New("empty")
	ErrNoInstance   = errors.New("no Instance Found")
	ErrSegmentFault = errors.New("segment Fault")
	ErrNoMemory     = errors.New("no memory")
	ErrProcessNil   = errors.New("process requested is nil")
	ErrNoTabla      = errors.New("no Tabla Found")
	ErrNoIndices    = errors.New("no Indices Found")
	ErrDuplicatePID = errors.New("asked to instance an already existing PID")
)
