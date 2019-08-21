package workerpool

import (
	"errors"
)

var (
	// ErrQueueAlreadyExist queue already exist error
	ErrQueueAlreadyExist = errors.New("queue already existed")
	// ErrQueueNotExist queue not exist error
	ErrQueueNotExist = errors.New("queue not exists")
)
