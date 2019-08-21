// Code generated by MockGen. DO NOT EDIT.
// Source: abstraction.go

// Package workerpool is a generated GoMock package.
package workerpool

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDelegator is a mock of IDelegator interface
type MockDelegator struct {
	ctrl     *gomock.Controller
	recorder *MockDelegatorMockRecorder
}

// MockDelegatorMockRecorder is the mock recorder for MockDelegator
type MockDelegatorMockRecorder struct {
	mock *MockDelegator
}

// NewMockDelegator creates a new mock instance
func NewMockDelegator(ctrl *gomock.Controller) *MockDelegator {
	mock := &MockDelegator{ctrl: ctrl}
	mock.recorder = &MockDelegatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDelegator) EXPECT() *MockDelegatorMockRecorder {
	return m.recorder
}

// Delegate mocks base method
func (m *MockDelegator) Delegate(ctx context.Context, queueName string, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delegate", ctx, queueName, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delegate indicates an expected call of Delegate
func (mr *MockDelegatorMockRecorder) Delegate(ctx, queueName, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delegate", reflect.TypeOf((*MockDelegator)(nil).Delegate), ctx, queueName, data)
}

// MockWorkerPool is a mock of IWorkerPool interface
type MockWorkerPool struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerPoolMockRecorder
}

// MockWorkerPoolMockRecorder is the mock recorder for MockWorkerPool
type MockWorkerPoolMockRecorder struct {
	mock *MockWorkerPool
}

// NewMockWorkerPool creates a new mock instance
func NewMockWorkerPool(ctrl *gomock.Controller) *MockWorkerPool {
	mock := &MockWorkerPool{ctrl: ctrl}
	mock.recorder = &MockWorkerPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorkerPool) EXPECT() *MockWorkerPoolMockRecorder {
	return m.recorder
}

// AddWorker mocks base method
func (m *MockWorkerPool) AddWorker(ctx context.Context, queueName string, worker Worker) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWorker", ctx, queueName, worker)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddWorker indicates an expected call of AddWorker
func (mr *MockWorkerPoolMockRecorder) AddWorker(ctx, queueName, worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWorker", reflect.TypeOf((*MockWorkerPool)(nil).AddWorker), ctx, queueName, worker)
}

// MockQueue is a mock of IQueue interface
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Enqueue mocks base method
func (m *MockQueue) Enqueue(ctx context.Context, data []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue
func (mr *MockQueueMockRecorder) Enqueue(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockQueue)(nil).Enqueue), ctx, data)
}

// Dequeue mocks base method
func (m *MockQueue) Dequeue(ctx context.Context) (<-chan []byte, chan<- error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dequeue", ctx)
	ret0, _ := ret[0].(<-chan []byte)
	ret1, _ := ret[1].(chan<- error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Dequeue indicates an expected call of Dequeue
func (mr *MockQueueMockRecorder) Dequeue(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dequeue", reflect.TypeOf((*MockQueue)(nil).Dequeue), ctx)
}

// Dispose mocks base method
func (m *MockQueue) Dispose() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispose")
}

// Dispose indicates an expected call of Dispose
func (mr *MockQueueMockRecorder) Dispose() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispose", reflect.TypeOf((*MockQueue)(nil).Dispose))
}