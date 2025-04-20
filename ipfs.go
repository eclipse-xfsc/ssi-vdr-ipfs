package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func IpfsService(client *http.Client) *ipfsService {
	return &ipfsService{client: client}
}

type ipfsService struct {
	client *http.Client
}

func (s *ipfsService) GarbageCollection() error {
	url := GetRpcEndpointURL(CollectGarbageEndpoint)
	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil || resp.StatusCode >= http.StatusMultipleChoices {
		if err != nil {
			err = errors.Join(fmt.Errorf("error during garbage collection"), err)
		} else {
			b, _ := io.ReadAll(resp.Body)
			err = errors.Join(fmt.Errorf("error during garbage collection"), fmt.Errorf("response status %s with body %s", resp.Status, string(b)))
		}
		return err
	}
	return nil
}

func (s *ipfsService) Ping() error {
	url := GetRpcEndpointURL("/api/v0/version")
	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil || resp.StatusCode >= http.StatusMultipleChoices {
		if err != nil {
			err = fmt.Errorf("not alive: %v", err)
		} else {
			err = fmt.Errorf("not alive: %s", resp.Status)
		}
		return err
	}
	return nil
}
