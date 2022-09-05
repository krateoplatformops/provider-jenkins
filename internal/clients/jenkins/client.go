package jenkins

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/krateoplatformops/provider-jenkins/internal/helpers"
	httphelper "github.com/krateoplatformops/provider-jenkins/internal/helpers/http"
)

type Client struct {
	baseUrl    string
	node       string
	username   *string
	password   *string
	httpClient *http.Client
	crumbData  map[string]string
}

func (c *Client) GetJobConfig(ctx context.Context, name string) ([]byte, error) {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}
	uri.Path = path.Join(uri.Path, c.node, "job", name, "config.xml")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	if c.username != nil && c.password != nil {
		req.SetBasicAuth(helpers.StringValue(c.username), helpers.StringValue(c.password))
	}

	return c.send(req)
}

func (c *Client) CreateJob(ctx context.Context, name string, data []byte) error {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	uri.Path = path.Join(uri.Path, c.node, "createItem")

	query := make(url.Values)
	query.Set("name", name)
	uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(data))
	if err != nil {
		return err
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return httphelper.NewErr(rsp)
}

func (c *Client) DeleteJob(ctx context.Context, name string) error {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	uri.Path = path.Join(uri.Path, c.node, "job", name)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri.String(), nil)
	if err != nil {
		return err
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return httphelper.NewErr(rsp)
}

func (c *Client) send(req *http.Request) ([]byte, error) {
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = httphelper.NewErr(rsp); err != nil {
		return nil, err
	}

	if rsp.Body == nil {
		return nil, nil
	}
	defer rsp.Body.Close()

	return io.ReadAll(rsp.Body)
}
