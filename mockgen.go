package workerpool

//go:generate mockgen -source=./internal/abs/abstraction.go -destination=./internal/mock/mock.go -package mock -mock_names IQueue=MockQueue,IDelegator=MockDelegator,IWorkerPool=MockWorkerPool
