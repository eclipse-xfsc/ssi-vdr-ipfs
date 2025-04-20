package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eclipse-xfsc/ssi-vdr-core/types"
	pathUtil "github.com/ipfs/boxo/path"
	cidUtil "github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/kubo/core/coreiface/options"
	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var Plugin = plugin{} //export Plugin Symbol

type plugin struct{}

func (p *plugin) GetVerifiableDataRegistry() (types.VerifiableDataRegistry, error) {
	// load env config and define logger
	setConfig()

	ipfs := &IPFSVerifiableDataRegistry{
		Options: DefaultIpfsConfig.Options,
		AddFile: DefaultIpfsConfig.AddFile,
		Client:  DefaultIpfsConfig.Client,
	}
	ipfs, err := ipfs.withShell()
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return ipfs, nil
}

var DefaultIpfsConfig = &IPFSVerifiableDataRegistryConfig{
	Options: []options.ApiOption{},
	AddFile: DefaultAddFile,
	Client:  http.DefaultClient,
}

func (r *IPFSVerifiableDataRegistry) withShell() (*IPFSVerifiableDataRegistry, error) {
	host := viper.Get("IPFS_HOST")
	port := viper.Get("IPFS_RPC_API_PORT")
	addr := fmt.Sprintf(MultiAddressTemplate, host, port)
	multiAddress, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Errorf("Error creating multiaddress %s", err.Error())
		return nil, err
	}
	node, err := rpc.NewApiWithClient(multiAddress, r.Client)
	if err != nil {
		return nil, err
	}
	nodeWithOptions, _ := node.WithOptions(r.Options...)
	r.Shell = nodeWithOptions.(*rpc.HttpApi)
	return r, nil
}

func (r *IPFSVerifiableDataRegistry) Get(cid *types.DataIdentifier) (*types.VDROutput, error) {
	url := fmt.Sprintf("%s/%s", viper.Get("IPFS_API_GATEWAY_URL"), cid.Value)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pinned, err := r.Shell.Pin().Ls(ctx)
	if err != nil {
		err = errors.Join(fmt.Errorf("Failed to list pinned"), err)
		return nil, err
	}
	isExists := IncludesCid(pinned, cid)
	if !isExists {
		return nil, types.DataIdentifierNotFound
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Failed to Get %s. %s", cid, err.Error())
		return nil, err
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		log.Errorf("Failed to Get %s. %s", cid, err.Error())
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to Get %s. %s", cid, err.Error())
		return nil, err
	}
	log.Debugf("Received file of length %v", len(data))
	out := &types.VDROutput{Data: data}
	return out, nil
}

func IncludesCid(pinned <-chan iface.Pin, cid *types.DataIdentifier) bool {
	for pin := range pinned {
		if pin.Path().RootCid().String() == cid.Value {
			return true
		}
	}
	return false
}

func (r *IPFSVerifiableDataRegistry) Put(cid *types.DataIdentifier, file io.Reader) (*types.DataIdentifier, error) {
	return r.AddFile(r.Shell, file, AddFsOptionPin(false), AddPinOptionRecursive(true))
}

func (r *IPFSVerifiableDataRegistry) Update(cid *types.DataIdentifier, file io.Reader) (*types.DataIdentifier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	oldCidObj, err := cidUtil.Decode(cid.Value)
	if err != nil {
		log.Errorf("Provided cid <%s> should reprent a valid CID", cid.Value)
		return nil, err
	}

	newPath, err := AddFileToFs(r.Shell, file, AddFsOptionPin(false), ctx)
	if err != nil {
		log.Errorf("Could not add new version of file with old %s: %s", cid.Format, cid.Value)
		return nil, err
	}

	pinner := r.Shell.Pin()
	oldPath := pathUtil.FromCid(oldCidObj)
	id := &types.DataIdentifier{
		Format: IdentifierFormatCID,
		Value:  newPath.RootCid().String(),
	}
	return id, pinner.Update(ctx, oldPath, newPath)
}

func (r *IPFSVerifiableDataRegistry) Delete(cid *types.DataIdentifier) error {
	pinner := r.Shell.Pin()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cidObj, err := cidUtil.Decode(cid.Value)
	if err != nil {
		err = errors.Join(fmt.Errorf("provided cid <%s> should reprent a valid CID", cid.Value), err)
		log.Error(err)
		return err
	}
	path := pathUtil.FromCid(cidObj)
	err = pinner.Rm(ctx, path)
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to unpin CID %s", cid.Value), err)
		log.Error(err)
		return err
	}
	err = IpfsService(r.Client).GarbageCollection()
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to collect garbage after unpin CID %s", cid.Value), err)
		log.Error(err)
		return err
	}
	return nil
}

func (r *IPFSVerifiableDataRegistry) IsAlive() bool {
	err := IpfsService(r.Client).Ping()
	return err == nil
}

func (r *IPFSVerifiableDataRegistry) Configure(newConfig interface{}) error {
	conf, ok := newConfig.(*IPFSVerifiableDataRegistryConfig)
	if !ok {
		return fmt.Errorf("could not convert passed newConfig to IPFSVerifiableDataRegistryConfig")
	}
	config := r.getIpfsConfig(conf)
	r.Options = config.Options
	r.AddFile = config.AddFile
	r.Client = config.Client
	_, err := r.withShell()
	return err
}

func (r *IPFSVerifiableDataRegistry) List() ([]*types.DataIdentifier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	items, err := r.Shell.Pin().Ls(ctx)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to list documents"), err)
	}
	res := make([]*types.DataIdentifier, 0)
	for item := range items {
		res = append(res, &types.DataIdentifier{
			Format: IdentifierFormatCID,
			Value:  item.Path().RootCid().String(),
		})
	}
	return res, nil
}

func (r *IPFSVerifiableDataRegistry) getIpfsConfig(conf *IPFSVerifiableDataRegistryConfig) *IPFSVerifiableDataRegistryConfig {
	var res = &IPFSVerifiableDataRegistryConfig{}
	if conf.AddFile == nil {
		res.AddFile = r.AddFile
	} else {
		res.AddFile = conf.AddFile
	}
	if conf.Client == nil {
		res.Client = r.Client
	} else {
		res.Client = conf.Client
	}
	if conf.Options == nil {
		res.Options = r.Options
	} else {
		res.Options = conf.Options
	}
	return res
}
