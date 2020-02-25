package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Direction int32

const (
	FORWARD  Direction = 0
	BACKWARD Direction = 1
)

var (
	Direction_value = map[string]int32{
		"FORWARD":  0,
		"BACKWARD": 1,
	}
	Direction_name = map[int32]string{
		0: "FORWARD",
		1: "BACKWARD",
	}
)

func (x Direction) String() string {
	s, ok := Direction_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

type queryrange struct {
	query      string
	start, end time.Time
	direction  Direction
	limit      uint32
	//
	step time.Duration
	url  url.URL
}

func doQueryRange(ctx context.Context, r queryrange, c *http.Client) (string, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", r.url.String(), http.NoBody)
	if err != nil {
		return "", nil
	}
	params := req.URL.Query()
	params.Add("start", fmt.Sprintf("%d", r.start.UnixNano()))
	params.Add("end", fmt.Sprintf("%d", r.end.UnixNano()))
	params.Add("query", r.query)
	params.Add("direction", r.direction.String())
	params.Add("limit", fmt.Sprintf("%d", r.limit))

	if r.step != 0 {
		params.Add("step", fmt.Sprintf("%f", r.step.Seconds()))
	}
	req.URL.Path = "/loki/api/v1/query_range"
	req.URL.RawQuery = params.Encode()
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode/2 != 100 {
		return string(bodyBytes), fmt.Errorf("status code fail: %d", resp.StatusCode)
	}

	return string(bodyBytes), nil
}

type labels struct {
	start, end time.Time
	name       string
	url        url.URL
}

func doLabels(ctx context.Context, r labels, c *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.url.String(), http.NoBody)
	if err != nil {
		return "", nil
	}
	params := req.URL.Query()
	params.Add("start", fmt.Sprintf("%d", r.start.UnixNano()))
	params.Add("end", fmt.Sprintf("%d", r.end.UnixNano()))

	req.URL.Path = "/loki/api/v1/label"
	if r.name != "" {
		req.URL.Path = fmt.Sprintf("/loki/api/v1/label/%s/values", r.name)
	}

	req.URL.RawQuery = params.Encode()
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode/2 != 100 {
		return string(bodyBytes), fmt.Errorf("status code fail: %d", resp.StatusCode)
	}

	return string(bodyBytes), nil
}
