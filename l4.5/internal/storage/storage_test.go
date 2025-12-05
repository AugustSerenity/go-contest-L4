package storage

import (
	"fmt"
	"testing"
	"time"

	"l2.18/internal/model"
)

func BenchmarkExactEventExists(b *testing.B) {
	st := New()
	now := time.Now()

	for i := 0; i < 100000; i++ {
		st.Create(1, model.Event{
			Name: fmt.Sprintf("ev_%d", i),
			Date: now.Add(time.Duration(i) * time.Hour),
		})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		st.ExactEventExists(1, now, "ev_1")
	}
}

func BenchmarkEventsForDay(b *testing.B) {
	st := New()
	now := time.Now()

	for i := 0; i < 100000; i++ {
		st.Create(1, model.Event{
			Name: fmt.Sprintf("ev_%d", i),
			Date: now.Add(time.Duration(i) * time.Hour),
		})
	}

	date := now.Add(48 * time.Hour) // 2-й день

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = st.EventsForDay(1, date)
	}
}
