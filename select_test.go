package selector

import "testing"

func BenchmarkSelectStmt(b *testing.B) {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	c4 := make(chan int)
	c5 := make(chan int)
	go func() {
		for {
			c3 <- 42
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-c1:
		case <-c2:
		case <-c3:
		case <-c4:
		case c5 <- 42:
		}
	}
}

func BenchmarkSelect(b *testing.B) {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	c4 := make(chan int)
	c5 := make(chan int)
	go func() {
		for {
			c3 <- 42
		}
	}()
	selector := New()
	selector.Add(c1, nil, nil)
	selector.Add(c2, nil, nil)
	selector.Add(c3, nil, nil)
	selector.Add(c4, nil, nil)
	selector.Add(c5, nil, func() int {
		return 42
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selector.Select()
	}
}
