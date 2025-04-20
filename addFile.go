package main

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/eclipse-xfsc/ssi-vdr-core/types"
	"github.com/ipfs/boxo/files"
	pathUtil "github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
)

type AddCommandResponse struct {
	CID  string `json:"Hash"`
	Name string
	Size string
}

type AddFunc func(shell *rpc.HttpApi, file io.Reader, addFsOptions options.UnixfsAddOption, addPinOptions options.PinAddOption) (*types.DataIdentifier, error)

var AddFsOptionPin = func(pin bool) options.UnixfsAddOption {
	return func(settings *options.UnixfsAddSettings) error {
		settings.Pin = pin
		return nil
	}
}

var AddPinOptionRecursive = func(recursive bool) options.PinAddOption {
	return func(settings *options.PinAddSettings) error {
		settings.Recursive = recursive
		return nil
	}
}

func DefaultAddFile(shell *rpc.HttpApi, file io.Reader, addFsOptions options.UnixfsAddOption, addPinOptions options.PinAddOption) (*types.DataIdentifier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	path, err := AddFileToFs(shell, file, addFsOptions, ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to add file to ipfs"), err)
	}
	err = shell.Pin().Add(ctx, path, addPinOptions)
	if err != nil {
		return nil, errors.Join(errors.New("failed to pin file to ipfs node"), err)
	}
	return &types.DataIdentifier{
		Format: IdentifierFormatCID,
		Value:  path.RootCid().String(),
	}, nil
}

func AddFileToFs(shell *rpc.HttpApi, file io.Reader, addFsOptions options.UnixfsAddOption, ctx context.Context) (pathUtil.ImmutablePath, error) {
	fileNode := files.NewReaderFile(file)
	path, err := shell.Unixfs().Add(ctx, fileNode, addFsOptions)
	return path, err
}
