package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithQueue(t *testing.T) {
	q := &fakeQueue{}
	opt := WithQueue("name", q)

	assert.NotNil(t, opt)
}