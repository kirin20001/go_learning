package concurrency

import "testing"

func TestPrintIota(t *testing.T) {
	printIota()
}

func TestTryLockMutex(t *testing.T) {
	tryLockMutex()
}

func TestCountWithLockInfo(t *testing.T) {
	countWithLockInfo()
}

func TestTryMutexQueue(t *testing.T) {
	tryMutexQueue()
}