package api

import (
	"errors"
	"fmt"
	"github.com/kylin-ops/grequests"
	"strings"
)

type Client struct {
	AccessAddress string
	AccessToken   string
}

const (
	levelGUEST      = 10
	levelREPORTER   = 20
	levelDEVELOPER  = 30
	levelMAINTAINER = 40
	levelOWNER      = 50

	visibilityPrivate  = "private"
	visibilityInternal = "internal"
	visibilityPublic   = "public"
)

func (c *Client) request(method, uri string, data ...grequests.Data) (*grequests.Response, error) {
	var resp *grequests.Response
	var err error
	var option = &grequests.RequestOptions{Header: grequests.Header{"PRIVATE-TOKEN": c.AccessToken}}
	var url = c.AccessAddress + uri
	if len(data) > 1 {
		option.Data = data[1]
	}
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = grequests.Get(url, option)
	case "POST":
		option.Form = true
		resp, err = grequests.Post(url, option)
	case "PUT":
		option.Form = true
		resp, err = grequests.Put(url, option)
	case "DELETE":
		resp, err = grequests.Delete(url, option)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	if resp.StatusCode() != 200 {
		if text, err := resp.Text(); err != nil {
			return resp, err
		} else {
			return resp, errors.New(text)
		}
	}
	return resp, nil
}

func (c *Client) getHeader() (header grequests.Header) {
	header = map[string]string{"PRIVATE-TOKEN": c.AccessToken}
	return header
}

func (c *Client) ProjectList() error {
	req, err := c.request("get", "/api/v4/projects")
	if err != nil {
		return err
	}
	fmt.Println(req.Text())
	return err
}

func (c *Client) ProjectGet(projectId string) error {
	req, err := c.request("get", "/api/v4/projects/"+projectId)
	if err != nil {
		return err
	}
	fmt.Println(req.Text())
	return err
}

func (c *Client) ProjectCreate(projectName string) error {
	var data = grequests.Data{"name": projectName, "path": projectName, "visibility": visibilityPrivate}
	req, err := c.request("post", "/api/v4/projects", data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(req.Text())
	return err
}
