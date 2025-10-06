package cache

import "github.com/valkey-io/valkey-go"

type ValkeyProvider interface {
	Client() (valkey.Client, error)
	ResetClient()
}
