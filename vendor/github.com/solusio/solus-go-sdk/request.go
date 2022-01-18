package solus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var maxRetries = 10

var errMaxRetriesReached = errors.New("exceeded retry limit")

// Func represents functions that can be retried.
type retryFunc func(attempt int) (retry bool, err error)

type requestOpts struct {
	params map[string][]string
	body   interface{}
}

type requestOption func(*requestOpts)

func withFilter(f map[string]string) requestOption {
	return func(o *requestOpts) {
		if f == nil {
			return
		}

		if o.params == nil {
			o.params = map[string][]string{}
		}
		for field, value := range f {
			o.params[field] = append(o.params[field], value)
		}
	}
}

func withBody(b interface{}) requestOption {
	return func(o *requestOpts) {
		o.body = b
	}
}

func (c *Client) create(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusCreated {
		return newHTTPError(http.MethodPost, path, code, body)
	}

	return unmarshal(body, &resp)
}

func (c *Client) list(ctx context.Context, path string, resp interface{}, opts ...requestOption) error {
	body, code, err := c.request(ctx, http.MethodGet, path, opts...)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodGet, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) get(ctx context.Context, path string, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodGet, path)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodGet, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) update(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPut, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodPut, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) patch(ctx context.Context, path string, data, resp interface{}) error {
	body, code, err := c.request(ctx, http.MethodPatch, path, withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(http.MethodPatch, path, code, body)
	}

	return unmarshal(body, resp)
}

func (c *Client) asyncDelete(ctx context.Context, path string) (Task, error) {
	body, code, err := c.request(ctx, http.MethodDelete, path)
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(http.MethodDelete, path, code, body)
	}

	var resp taskResponse
	if err := unmarshal(body, &resp); err != nil {
		return Task{}, err
	}

	if resp.Data.ID == 0 {
		return Task{}, errors.New("task doesn't have an id")
	}

	return resp.Data, nil
}

func (c *Client) syncDelete(ctx context.Context, path string) error {
	body, code, err := c.request(ctx, http.MethodDelete, path)
	if err != nil {
		return err
	}

	if code != http.StatusNoContent {
		return newHTTPError(http.MethodDelete, path, code, body)
	}
	return nil
}

func (c *Client) asyncPost(ctx context.Context, path string, opts ...requestOption) (Task, error) {
	body, code, err := c.request(ctx, http.MethodPost, path, opts...)
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp taskResponse
	if err := unmarshal(body, &resp); err != nil {
		return Task{}, err
	}

	if resp.Data.ID == 0 {
		return Task{}, errors.New("task doesn't have an id")
	}

	return resp.Data, nil
}

func (c *Client) request(ctx context.Context, method, path string, opts ...requestOption) ([]byte, int, error) {
	req, err := c.buildRequest(ctx, method, path, opts...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	var resp *http.Response
	err = retry(func(attempt int) (bool, error) {
		var err error
		resp, err = c.HTTPClient.Do(req)

		checkOk, checkErr := checkForRetry(resp, err)

		if checkErr != nil {
			err = checkErr
		}

		if checkOk {
			time.Sleep(c.RetryAfter)        // wait before next try
			return attempt < c.Retries, err // try N times
		}

		return false, err
	})
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Errorf("failed to close response body for %s %s: %s", method, path, err)
		}
	}()

	code := resp.StatusCode

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, code, nil
}

func (c *Client) buildRequest(ctx context.Context, method, path string, opts ...requestOption) (*http.Request, error) {
	reqOpts := requestOpts{}
	for _, o := range opts {
		o(&reqOpts)
	}

	var (
		bodyByte []byte
		reqBody  io.ReadWriter
	)
	if reqOpts.body != nil {
		var err error
		bodyByte, err = json.Marshal(reqOpts.body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(bodyByte)
	}

	url, err := c.buildURL(path, reqOpts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	for k, values := range c.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("User-Agent", c.UserAgent)

	c.Logger.Debugf("[%s] %s with body %q", method, url, string(bodyByte))
	return req, nil
}

func (c *Client) buildURL(path string, opts requestOpts) (string, error) {
	fullURL, err := c.BaseURL.Parse(path)
	if err != nil {
		return "", err
	}

	if opts.params != nil {
		query := fullURL.Query()
		for param, values := range opts.params {
			for _, value := range values {
				query.Add(param, value)
			}
		}

		fullURL.RawQuery = query.Encode()
	}
	return fullURL.String(), nil
}

func checkForRetry(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return false, nil
}

func retry(fn retryFunc) error {
	var (
		err   error
		retry bool
	)
	attempt := 1
	for {
		retry, err = fn(attempt)
		if !retry {
			break
		}
		attempt++
		if attempt > maxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

func unmarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to decode %q: %w", data, err)
	}
	return nil
}
