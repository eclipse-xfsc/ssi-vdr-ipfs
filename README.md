# VDR IPFS

## Description

A plugin implementing ```VerifiableDataRegistry``` interface using Interplanetary Filesystem for storing files. Used by [VDR](https://gitlab.eclipse.org/eclipse/xfsc/libraries/crypto/engine/core/-/blob/main/README.md?ref_type=heads)

## Building

Here is the [README.md](https://gitlab.eclipse.org/eclipse/xfsc/dev-ops/building/go-plugin/-/blob/main/README.md#building-go-services-with-plugin-based-dependencies) describing the specifics of build process for services, where the dependency is used.

## Dependencies

The package should establish connection with [IPFS node](https://ipfs.tech)

## Entities

Exports ```GetVerifiableDataRegistry()``` that returns an instance of ```VerifiableDataRegistry```. 
Used for CRUD operations in IPFS.

