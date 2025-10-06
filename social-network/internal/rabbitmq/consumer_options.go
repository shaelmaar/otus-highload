package rabbitmq

type ConsumerOption[T Message] func(*Consumer[T])

func WithPrefetchCount[T Message](prefetchCount int) ConsumerOption[T] {
	return func(c *Consumer[T]) {
		c.prefetchCount = prefetchCount
	}
}

func WithWorkerCount[T Message](count int) ConsumerOption[T] {
	return func(c *Consumer[T]) {
		c.workerCount = count
	}
}
