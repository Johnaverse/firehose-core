package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
)

var ErrorNoMoreClient = errors.New("no more clients")

type Clients[C any] struct {
	clients               []C
	next                  int
	maxBlockFetchDuration time.Duration
}

func NewClients[C any](maxBlockFetchDuration time.Duration) *Clients[C] {
	return &Clients[C]{
		next:                  0,
		maxBlockFetchDuration: maxBlockFetchDuration,
	}
}

func (c *Clients[C]) Add(client C) {
	c.clients = append(c.clients, client)
}

func (c *Clients[C]) Next() (client C, err error) {
	if len(c.clients) <= c.next {
		return client, ErrorNoMoreClient
	}
	client = c.clients[c.next]
	c.next++
	return client, nil
}

func WithClients[C any, V any](clients *Clients[C], f func(context.Context, C) (v V, err error)) (v V, err error) {
	clients.next = 0
	var errs error
	for {
		client, err := clients.Next()
		if err != nil {
			errs = multierror.Append(errs, err)
			return v, errs
		}
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, clients.maxBlockFetchDuration)
		v, err := f(ctx, client)
		cancel()
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		return v, nil
	}
}
