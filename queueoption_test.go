package workerpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithWorker(t *testing.T) {
	const qName = "test"
	queueOption := WithWorker(context.TODO(), fakeWorker, 1)
	assert.NotNil(t, queueOption)
}

func fakeWorker(_ context.Context, _ []byte) error {
	return nil
}
