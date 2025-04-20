package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/ipfs/kubo/client/rpc"
)

var mockAddResponse = AddCommandResponse{CID: "testid", Name: "test", Size: "2"}

type requestBuilderMock struct{}

func (r *requestBuilderMock) Arguments(args ...string) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) BodyString(body string) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) BodyBytes(body []byte) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) Body(body io.Reader) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) FileBody(body io.Reader) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) Option(key string, value interface{}) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) Header(name, value string) rpc.RequestBuilder {
	return nil
}
func (r *requestBuilderMock) Send(ctx context.Context) (*rpc.Response, error) {
	return nil, nil
}
func (r *requestBuilderMock) Exec(ctx context.Context, obj interface{}) error {
	data := []byte(`{"Hash":"testcid"}`)
	return json.Unmarshal(data, obj)
}

func TestConstructRequest(t *testing.T) {
	defer monkey.UnpatchAll()
	var requestMock = func(c *rpc.HttpApi, command string, args ...string) rpc.RequestBuilder {
		if !(command == AddCommand && reflect.DeepEqual(args, []string{})) {
			t.Errorf("Expected to receive %s command and `[]` args. Actual command %s args %v", AddCommand, command, args)
		}
		return &requestBuilderMock{}
	}
	monkey.PatchInstanceMethod(reflect.TypeOf(&rpc.HttpApi{}), "Request", requestMock)
	shell, _ := rpc.NewURLApiWithClient("", http.DefaultClient)
	fileData := []byte("test file")
	file := bytes.NewReader(fileData)
	_, err := DefaultAddFile(shell, file, AddFsOptionPin(false), AddPinOptionRecursive(false))
	if err != nil {
		t.Error(err)
	}
}
