package workerpool

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueueItem(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
	}{
		{"not_nil", []byte("test")},
		{"nil", nil},
		{"empty", []byte{}},
	}

	for i := 0; i < len(testcases); i++ {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			qItem := NewQueueItem(tc.data)
			assert.NotNil(t, qItem)
			assert.Equal(t, tc.data, qItem.data)
		})
	}
}

func TestQueueItem_Data(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
	}{
		{"not_nil", []byte("test")},
		{"nil", nil},
		{"empty", []byte{}},
	}

	for i := 0; i < len(testcases); i++ {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			qItem := QueueItem{
				data: tc.data,
			}
			assert.NotNil(t, qItem)
			assert.Equal(t, tc.data, qItem.Data())
		})
	}
}

func TestQueueItem_Complete(t *testing.T) {
	testcases := []struct {
		name           string
		expectCalled   bool
		err error
		mockOnComplete *mockOnComplete
	}{
		{"not_nil", true, nil, &mockOnComplete{}},
		{"nil", true, errors.New("test"), &mockOnComplete{}},
	}

	for i := 0; i < len(testcases); i++ {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			qItem := QueueItem{	complete: tc.mockOnComplete.Complete}
			qItem.Complete(context.TODO(), tc.err)
			assert.Equal(t, tc.err, tc.mockOnComplete.err)
			assert.Equal(t, tc.expectCalled, tc.mockOnComplete.called)
		})
	}
}

func TestQueueItem_CompleteNil(t *testing.T) {
	qItem := QueueItem{	complete: nil}
	qItem.Complete(context.TODO(), nil)
	assert.True(t, true, "should not panic when complete nil")
}
func TestQueueItem_WhenComplete(t *testing.T) {
	qItem := QueueItem{	complete: nil}
	called := false
	fake := func(ctx context.Context, err error) {
		called = true
	}
	qItem.WhenComplete(fake)
	assert.NotNil(t, qItem.complete)

	qItem.complete(context.TODO(), nil)
	assert.True(t, called)
}

type mockOnComplete struct {
	err    error
	called bool
}

func (t *mockOnComplete) Complete(ctx context.Context, err error) {
		t.called = true
		t.err = err
}
