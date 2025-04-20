package main

// ipfs cli

const (
	AddCommand             = "add"
	CollectGarbageEndpoint = "/api/v0/repo/gc"
	RecursiveOption        = "--recursive"
	StdinNameOption        = "--stdin-name=%s"
	PinOption              = "--pin=%t"
)

const (
	EnvFilePath          = ".env"
	MultiAddressTemplate = "/dns4/%s/tcp/%s/http/"
	IdentifierFormatCID  = "cid"
)
