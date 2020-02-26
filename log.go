package main

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	tm "github.com/buger/goterm"
)

type logger struct {
	rt      http.RoundTripper
	counter *int64
}

func (l *logger) RoundTrip(r *http.Request) (*http.Response, error) {
	defer func() {
		_ = atomic.AddInt64(l.counter, 1)
	}()
	return l.rt.RoundTrip(r)
}

func Logger(ctx context.Context, t *time.Ticker, rt http.RoundTripper) http.RoundTripper {
	if !*verboseFlag {
		return rt
	}
	l := &logger{
		rt:      rt,
		counter: new(int64),
	}
	go func() {
		start := time.Now()
		last := int64(0)
		tm.Clear()
		for {
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					new := atomic.SwapInt64(l.counter, *l.counter)
					since := new - last
					last = new
					if since != 0 {
						tm.MoveCursor(1, 1)
						tm.Printf("Current RPS: %f  Total requests: %d", float64(since)/time.Since(start).Seconds(), new)
						tm.Flush()
					}
					start = time.Now()
				}
			}
		}
	}()
	return l
}
