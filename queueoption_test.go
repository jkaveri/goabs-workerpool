package workerpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithWorker(t *testing.T) {
	queueOption := WithWorker(context.TODO(), fakeWorker)
	assert.NotNil(t, queueOption)
}

func TestWithWorkers(t *testing.T) {
	opt := WithWorkers(context.TODO(), fakeWorker, 1)
	assert.NotNil(t, opt)
}

func fakeWorker(ctx context.Context, data []byte) error {
	return nil
}