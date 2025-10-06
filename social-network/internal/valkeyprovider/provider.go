package valkeyprovider

import (
	"errors"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

var ErrClientInactive = errors.New("client is inactive")

type Provider struct {
	clientOption valkey.ClientOption
	logger       *zap.Logger

	mx         sync.RWMutex
	once       sync.Once
	stopRetry  chan struct{}
	isActive   bool
	isRetrying bool
	client     valkey.Client
}

func NewProvider(
	clientOption valkey.ClientOption,
	logger *zap.Logger,
) (*Provider, error) {
	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	p := &Provider{
		clientOption: clientOption,
		logger:       logger,
		mx:           sync.RWMutex{},
		once:         sync.Once{},
		stopRetry:    make(chan struct{}),
		isActive:     false,
		isRetrying:   false,
		client:       nil,
	}

	return p, nil
}

func (p *Provider) Init() {
	p.initClient()
}

func (p *Provider) ResetClient() {
	p.mx.Lock()

	if p.isRetrying {
		p.mx.Unlock()

		return
	}

	if p.client != nil {
		p.client.Close()
		p.client = nil
	}

	p.isActive = false

	p.mx.Unlock()

	go p.retryConnect()
}

func (p *Provider) Client() (valkey.Client, error) {
	p.mx.RLock()
	defer p.mx.RUnlock()

	if p.isActive {
		return p.client, nil
	}

	return nil, ErrClientInactive
}

func (p *Provider) Close() {
	p.mx.Lock()
	defer p.mx.Unlock()

	close(p.stopRetry)

	if p.isActive {
		p.client.Close()
	}
}

func (p *Provider) initClient() {
	p.once.Do(func() {
		client, err := valkey.NewClient(p.clientOption)
		if err != nil {
			p.logger.Error("failed to init valkey client", zap.Error(err))

			go p.retryConnect()
		} else {
			p.mx.Lock()
			p.client = client
			p.isActive = true
			p.mx.Unlock()
		}
	})
}

func (p *Provider) retryConnect() {
	p.mx.Lock()
	p.isRetrying = true
	p.mx.Unlock()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopRetry:
			return
		case <-ticker.C:
			client, err := valkey.NewClient(p.clientOption)
			if err != nil {
				p.logger.Error("failed to init valkey client", zap.Error(err))
			} else {
				p.mx.Lock()
				p.client = client
				p.isActive = true
				p.isRetrying = false
				p.mx.Unlock()

				return
			}
		}
	}
}
