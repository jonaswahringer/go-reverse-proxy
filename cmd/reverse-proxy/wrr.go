package main

import (
	"errors"
	"fmt"
	"net/url"
	"sync/atomic"
)

type WRoundRobin interface {
	Next() *url.URL
}

type wroundrobin struct {
	urls []*url.URL
	next uint32
}

func New(urls ...*url.URL) (WRoundRobin, error) {
	if len(urls) == 0 {
		fmt.Println("error: no urls specified")
		return nil, errors.New("no servers available")
	}

	return &wroundrobin{
		urls: urls,
	}, nil

}

func (wrr *wroundrobin) Next() *url.URL {
	n := atomic.AddUint32(&wrr.next, 1)
	return wrr.urls[(int(n)-1)%len(wrr.urls)]
}
