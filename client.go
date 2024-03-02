package rpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

type ServiceClient struct {
	client    *http.Client
	endpoints map[string]string
}

func NewServiceClient(name, addr string) (*ServiceClient, error) {
	client := &ServiceClient{
		client: &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(
					network string,
					addr string,
					cfg *tls.Config,
				) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		},
	}

	endpoint, _ := url.JoinPath(addr, name)
	req, _ := http.NewRequest("POST", endpoint, nil)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to service: %v", err)
	}

	endpoints, err := DecodeJSON[map[string]string](resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range endpoints {
		endpoint, _ := url.JoinPath(addr, v)
		endpoints[k] = endpoint
	}

	client.endpoints = endpoints

	return client, nil
}

func (c *ServiceClient) Call(
	ctx context.Context,
	name string,
	payload io.Reader,
) (*http.Response, error) {
	endpoint, ok := c.endpoints[name]
	if !ok {
		return nil, errors.New("endpoint not found")
	}

	req, _ := http.NewRequest("POST", endpoint, payload)
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)
	return c.client.Do(req)
}

func Call[Req any, Res any](
	client *ServiceClient,
	req *Request[Req, Res],
) (*Res, error) {
	ctx := req.ctx
	name := req.name
	payload := req.payload

	reader := bytes.Buffer{}

	err := json.NewEncoder(&reader).Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %v", err)
	}

	resp, err := client.Call(ctx, name, &reader)
	if err != nil {
		return nil, fmt.Errorf("error calling procedure: %v", err)
	}

	res, err := DecodeJSON[response[Res]](resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	return res.Payload, nil
}
