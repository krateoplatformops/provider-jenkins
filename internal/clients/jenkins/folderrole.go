package jenkins

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type AddFolderRoleOpts struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions,omitempty"`
	FolderNames []string `json:"folderNames,omitempty"`
}

func (c *Client) AddFolderRole(ctx context.Context, opts AddFolderRoleOpts) error {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	uri.Path = path.Join(uri.Path, c.controller, "/folder-auth/addFolderRole")

	data, err := json.Marshal(&opts)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = c.doRequest(req)
	return err
}

func (c *Client) AssignSidToFolderRole(ctx context.Context, sid, roleName string) error {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	uri.Path = path.Join(uri.Path, c.controller, "/folder-auth/assignSidToFolderRole")

	params := url.Values{}
	params.Add("roleName", roleName)
	params.Add("sid", sid)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.doRequest(req)
	return err
}

func (c *Client) DeleteFolderRole(ctx context.Context, roleName string) error {
	uri, err := url.Parse(c.baseUrl)
	if err != nil {
		return err
	}
	uri.Path = path.Join(uri.Path, c.controller, "/folder-auth/deleteFolderRole")

	query := make(url.Values)
	query.Set("roleName", roleName)
	//uri.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), strings.NewReader(query.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.doRequest(req)
	return err
}
