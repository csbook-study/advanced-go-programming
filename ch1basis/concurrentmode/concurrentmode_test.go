package concurrentmode

import "testing"

func TestProducerConsumerBasis(t *testing.T) {
	t.Run("ProducerConsumerBasis", func(t *testing.T) {
		ProducerConsumerBasis()
	})
}

func TestPublishSubscribeBasis(t *testing.T) {
	t.Run("PublishSubscribeBasis", func(t *testing.T) {
		PublishSubscribeBasis()
	})
}

func TestPrimeSieveBasis(t *testing.T) {
	t.Run("PrimeSieveBasis", func(t *testing.T) {
		PrimeSieveBasis()
	})
}

func TestPrimeSieveContextBasis(t *testing.T) {
	t.Run("PrimeSieveContextBasis", func(t *testing.T) {
		PrimeSieveContextBasis()
	})
}
