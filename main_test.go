package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/eclipse-xfsc/ssi-vdr-core/types"
	pathUtil "github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
)

var (
	testFileData = "test data"
	testDataId   = &types.DataIdentifier{Format: IdentifierFormatCID, Value: "QmVLwvmGehsrNEvhcCnnsw5RQNseohgEkFNN1848zNzdng"}
	invalidCid   = &types.DataIdentifier{Format: IdentifierFormatCID, Value: "invalid cid"}
)

func TestGetVerifiableDataRegistryGet(t *testing.T) {
	var getMock = func(c *http.Client, url string) (*http.Response, error) {
		r := http.Response{
			Body: io.NopCloser(bytes.NewBufferString(testFileData)),
		}
		return &r, nil
	}
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", getMock)
	plgn := plugin{}
	registry, err := plgn.GetVerifiableDataRegistry()
	if err != nil {
		t.Error(err)
	}
	actualGet, _ := registry.Get(testDataId)
	if string(actualGet.Data) != testFileData {
		t.Errorf("Expected: %s actual: %s", testFileData, string(actualGet.Data))
	}
}

func TestGetVerifiableDataRegistryConfigureAndPut(t *testing.T) {
	plgn := plugin{}
	registry, err := plgn.GetVerifiableDataRegistry()
	if err != nil {
		t.Error(err)
	}
	var addFunc = func(shell *rpc.HttpApi, file io.Reader, addFsOptions options.UnixfsAddOption, addPinOptions options.PinAddOption) (*types.DataIdentifier, error) {
		return testDataId, nil
	}
	config := &IPFSVerifiableDataRegistryConfig{AddFile: addFunc}
	err = registry.Configure(config)
	if err != nil {
		t.Error(err)
	}
	actualPut, _ := registry.Put(testDataId, bytes.NewReader([]byte(testFileData)))
	if actualPut.Value != testDataId.Value {
		t.Errorf("Expected: %s actual: %s", testDataId, actualPut.Value)

	}
}

func TestGetVerifiableDataRegistryUpdate(t *testing.T) {
	var updateMock = func(p *rpc.PinAPI, ctx context.Context, from pathUtil.Path, to pathUtil.Path, opts ...options.PinUpdateOption) error {
		return nil
	}
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(&rpc.PinAPI{}), "Update", updateMock)

	plgn := plugin{}
	registry, err := plgn.GetVerifiableDataRegistry()
	if err != nil {
		t.Error(err)
	}
	var addFuncInvalidId = func(shell *rpc.HttpApi, file io.Reader, addFsOptions options.UnixfsAddOption, addPinOptions options.PinAddOption) (*types.DataIdentifier, error) {
		return testDataId, nil
	}
	config := &IPFSVerifiableDataRegistryConfig{AddFile: addFuncInvalidId}
	err = registry.Configure(config)
	if err != nil {
		t.Error(err)
	}
	_, err = registry.Update(invalidCid, bytes.NewReader([]byte(testFileData)))
	if err == nil {
		t.Errorf("Provided invalid cid %v. Error expected", invalidCid)
	}
	actual, err := registry.Update(testDataId, bytes.NewReader([]byte(testFileData)))
	if actual != testDataId {
		t.Errorf("Expected %v actual %v", testDataId, actual)
	}
}

func TestGetVerifiableDataRegistryDelete(t *testing.T) {
	var rmMock = func(p *rpc.PinAPI, ctx context.Context, path pathUtil.Path, opts ...options.PinRmOption) error {
		return nil
	}
	defer monkey.UnpatchAll()
	monkey.PatchInstanceMethod(reflect.TypeOf(&rpc.PinAPI{}), "Rm", rmMock)

	plgn := plugin{}
	registry, err := plgn.GetVerifiableDataRegistry()
	if err != nil {
		t.Error(err)
	}
	err = registry.Delete(invalidCid)
	if err == nil {
		t.Errorf("Provided invalid cid %v. Error expected", invalidCid)
	}
	err = registry.Delete(testDataId)
	if err != nil {
		t.Errorf("Expected nil actual %s", err.Error())
	}
}
