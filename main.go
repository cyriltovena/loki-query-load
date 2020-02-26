package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	urlFlag     = flag.String("url", "http://localhost:3100", "the url of the Loki server to target.")
	verboseFlag = flag.Bool("verbose", false, "print stats")
	client      = &http.Client{
		Timeout: 5 * time.Minute,
	}
)

func main() {
	flag.Parse()

	if urlFlag == nil {
		fmt.Println("-url is required.")
		os.Exit(1)
	}
	u, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Println("failed to parse url:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.Transport = Logger(ctx, time.NewTicker(time.Second), http.DefaultTransport)

	for i := 0; i < 12; i++ {
		go worker(ctx, *u)
	}
	<-ctx.Done()
}

func worker(ctx context.Context, u url.URL) {
	for {
		if ctx.Err() != nil {
			return
		}
		err := queryLabels(ctx, u)
		if err != nil {
			fmt.Println("err received: ", err)
		}
		err = query(ctx, u)
		if err != nil {
			fmt.Println("err received: ", err)
		}
	}

}

func query(ctx context.Context, u url.URL) error {
	_, err := doQueryRange(ctx, queryrange{
		start:     time.Now().Add(-6 * time.Hour),
		end:       time.Now(),
		direction: BACKWARD,
		limit:     10000,
		query:     `{namespace="cortex-ops"} |= "foo" != "foo"`,
		step:      1 * time.Minute,
		url:       u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doQueryRange(ctx, queryrange{
		start:     time.Now().Add(-24 * time.Hour),
		end:       time.Now(),
		direction: BACKWARD,
		limit:     10000,
		query:     `{namespace="default"} |= "foo" != "foo"`,
		url:       u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doQueryRange(ctx, queryrange{
		start:     time.Now().Add(-1 * time.Hour),
		end:       time.Now(),
		direction: BACKWARD,
		limit:     10000,
		query:     `{namespace="default"} |= "foo" != "foo"`,
		step:      5 * time.Second,
		url:       u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doQueryRange(ctx, queryrange{
		start:     time.Now().Add(-1 * time.Hour),
		end:       time.Now(),
		direction: BACKWARD,
		limit:     10000,
		query:     `{cluster="dev-us-central1", job=~".*/query-frontend"}`,
		step:      5 * time.Second,
		url:       u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doQueryRange(ctx, queryrange{
		start:     time.Now().Add(-1 * time.Hour),
		end:       time.Now(),
		direction: BACKWARD,
		limit:     10000,
		query:     `{job="default/nginx", namespace="default"}`,
		step:      5 * time.Second,
		url:       u,
	}, client)

	if err != nil {
		return err
	}

	return nil
}

func queryLabels(ctx context.Context, u url.URL) error {
	_, err := doLabels(ctx, labels{
		start: time.Now().Add(-1 * time.Hour),
		end:   time.Now(),
		url:   u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doLabels(ctx, labels{
		start: time.Now().Add(-1 * time.Hour),
		end:   time.Now(),
		name:  "status_code",
		url:   u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doLabels(ctx, labels{
		start: time.Now().Add(-3 * time.Hour),
		end:   time.Now(),
		name:  "container_name",
		url:   u,
	}, client)

	if err != nil {
		return err
	}

	_, err = doLabels(ctx, labels{
		start: time.Now().Add(-6 * time.Hour),
		end:   time.Now(),
		name:  "job",
		url:   u,
	}, client)

	if err != nil {
		return err
	}
	return nil
}
