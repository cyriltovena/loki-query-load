package main

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"
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

func LogRPS(ctx context.Context, t *time.Ticker, rt http.RoundTripper) http.RoundTripper {
	l := &logger{
		rt:      rt,
		counter: new(int64),
	}
	go func() {
		for {
			start := time.Now()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					log.Printf("Current RPS: %f", float64(*l.counter)/time.Since(start).Seconds())
					start = time.Now()
				}
			}
		}
	}()
	return l
}
