package dlock

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	l := New(":6379")

	id, err := l.Lock("k1", 100*time.Millisecond)
	if err != nil {
		t.Error(err)
	}
	defer l.Unlock("k1", id)

	if _, err := l.Lock("k1", 100*time.Millisecond); err == nil {
		t.Fatal("must be failed")
	}
}

func TestTimeout(t *testing.T) {
	l := New(":6379")

	if _, err := l.Lock("k2", 100*time.Millisecond); err != nil {
		t.Error(err)
	}

	time.Sleep(200 * time.Millisecond)

	id, err := l.Lock("k2", 100*time.Millisecond)
	if err != nil {
		t.Fatal("must be acquired")
	}
	l.Unlock("k2", id)
}
