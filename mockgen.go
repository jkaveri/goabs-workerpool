package workerpool

//go:generate mockgen -source=abstraction.go -destination=mock.go -package workerpool -mock_names IQueue=MockQueue,IDelegator=MockDelegator,IWorkerPool=MockWorkerPool
