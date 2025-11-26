package or

import (
	"testing"
	"time"
)

func helperChanClose(after time.Duration) <-chan interface{} {
	done := make(chan interface{})
	go func() {
		time.Sleep(after)
		close(done)
	}()
	return done
}

func TestOrSingleChannel(t *testing.T) {
	done := Or(helperChanClose(50 * time.Millisecond))

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Or didn't close for single channel")
	}
}

func TestOrMultipleChannels(t *testing.T) {
	start := time.Now()

	<-Or(
		helperChanClose(500*time.Millisecond),
		helperChanClose(80*time.Millisecond),
		helperChanClose(300*time.Millisecond),
	)

	check := time.Since(start)

	if check > 200*time.Millisecond {
		t.Fatalf("Or closed too slow, expected ~80ms, got %v", check)
	}
}

func TestOrNoChannels(t *testing.T) {
	done := Or()

	select {
	case <-done:
	default:
		t.Fatalf("Or() with no channels must close immediately")
	}
}
