package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/ipfs/kubo/client/rpc"
)

type Requester interface {
	Request(command string, args ...string) rpc.RequestBuilder
}

/*****************************
 * Test-Mocks
 *****************************/

// Mock für RequestBuilder
type requestBuilderMock struct{}

func (r *requestBuilderMock) Arguments(args ...string) rpc.RequestBuilder { return r }
func (r *requestBuilderMock) BodyString(body string) rpc.RequestBuilder   { return r }
func (r *requestBuilderMock) BodyBytes(body []byte) rpc.RequestBuilder    { return r }
func (r *requestBuilderMock) Body(body io.Reader) rpc.RequestBuilder      { return r }
func (r *requestBuilderMock) FileBody(body io.Reader) rpc.RequestBuilder  { return r }
func (r *requestBuilderMock) Option(key string, value interface{}) rpc.RequestBuilder {
	return r
}
func (r *requestBuilderMock) Header(name, value string) rpc.RequestBuilder { return r }
func (r *requestBuilderMock) Send(ctx context.Context) (*rpc.Response, error) {
	return nil, nil
}
func (r *requestBuilderMock) Exec(ctx context.Context, obj interface{}) error {
	// Liefert ein JSON, das in AddCommandResponse passt
	data := []byte(`{"Hash":"testcid","Name":"test","Size":"2"}`)
	return json.Unmarshal(data, obj)
}

// Mock für HttpApi (ersetzt rpc.HttpApi im Test)
type httpApiMock struct {
	*rpc.HttpApi
	lastCommand string
	lastArgs    []string
}

func (h *httpApiMock) Request(command string, args ...string) rpc.RequestBuilder {
	h.lastCommand = command
	h.lastArgs = args
	return &requestBuilderMock{}
}

/*****************************
 * Tests
 *****************************/

func TestConstructRequest(t *testing.T) {
	shell := &httpApiMock{}

	fileData := []byte("test file")
	file := bytes.NewReader(fileData)

	resp, err := DefaultAddFile(shell.HttpApi, file, AddFsOptionPin(false), AddPinOptionRecursive(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Command und Args validieren
	if shell.lastCommand != AddCommand || len(shell.lastArgs) != 0 {
		t.Errorf("Expected command=%s args=[], got command=%s args=%v", AddCommand, shell.lastCommand, shell.lastArgs)
	}

	// Response validieren
	if resp.Format != "testcid" {
		t.Errorf("Expected CID=testcid, got %s", resp.Format)
	}
	if resp.Value != "test" {
		t.Errorf("Expected Name=test, got %s", resp.Value)
	}

}
