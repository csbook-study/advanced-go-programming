package concurrentmemory

import "testing"

func TestMutexBasis(t *testing.T) {
	t.Run("MutexBasis", func(t *testing.T) {
		MutexBasis()
	})
}

func TestAtomicBasis(t *testing.T) {
	t.Run("AtomicBasis", func(t *testing.T) {
		AtomicBasis()
	})
}

func TestSingletonBasis(t *testing.T) {
	t.Run("SingletonBasis", func(t *testing.T) {
		SingletonBasis()
	})
}

func TestConfigBasis(t *testing.T) {
	t.Run("ConfigBasis", func(t *testing.T) {
		ConfigBasis()
	})
}

func TestChannelBasis(t *testing.T) {
	t.Run("ChannelBasis", func(t *testing.T) {
		ChannelBasis()
	})
}

func TestWaitGroupBasis(t *testing.T) {
	t.Run("WaitGroupBasis", func(t *testing.T) {
		WaitGroupBasis()
	})
}

func TestContextBasis(t *testing.T) {
	t.Run("ContextBasis", func(t *testing.T) {
		ContextBasis()
	})
}
