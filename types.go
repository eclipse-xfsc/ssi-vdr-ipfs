package main

import (
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	"net/http"
)

type IPFSVerifiableDataRegistry struct {
	Shell   *rpc.HttpApi
	Options []options.ApiOption
	AddFile AddFunc
	Client  *http.Client
}

type IPFSVerifiableDataRegistryConfig struct {
	Options []options.ApiOption
	AddFile AddFunc
	Client  *http.Client
}
