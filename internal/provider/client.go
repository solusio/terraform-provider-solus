// Copyright 1999-2024. WebPros International GmbH. All rights reserved.

package provider

import (
	"context"
	"net/url"
	"sync"

	"github.com/solusio/solus-go-sdk"
)

type client struct {
	*solus.Client

	account struct {
		once sync.Once
		err  error
		data solus.User
	}
}

func newClient(baseURL *url.URL, a solus.Authenticator, opts ...solus.ClientOption) (*client, error) {
	c, err := solus.NewClient(baseURL, a, opts...)
	if err != nil {
		return nil, err
	}

	return &client{
		Client: c,
	}, nil
}

func (c *client) CurrentUser(ctx context.Context) (solus.User, error) {
	c.account.once.Do(func() {
		c.account.data, c.account.err = c.Account.Get(ctx)
	})
	return c.account.data, c.account.err
}
