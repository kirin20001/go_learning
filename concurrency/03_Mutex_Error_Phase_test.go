package concurrency

import "testing"

func TestMutexError(t *testing.T) {
	mutexCounter()
}

func TestMutexMultiLock(t *testing.T) {
	mutexMultiLock()
}

func TestRecursiveMutex(t *testing.T) {
	recMutex()
}

func TestDeadLockScenario(t *testing.T) {
	deadLockScenario()
}