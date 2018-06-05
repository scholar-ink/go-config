package vault

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/micro/go-config/source"
)

// Currently a single vault reader
type vault struct {
	prefix      string
	stripPrefix string
	addr        string
	opts        source.Options
	client      *api.Client
}

var (
	DefaultPrefix = "/secret/micro/config/"
)

func (v *vault) Read() (*source.ChangeSet, error) {
	kv, err := v.client.Logical().List(v.prefix)
	if err != nil {
		return nil, err
	}

	if kv == nil || len(kv.Data) == 0 {
		return nil, fmt.Errorf("source not found: %s", v.prefix)
	}

	b, err := v.opts.Encoder.Encode(kv.Data)
	if err != nil {
		return nil, fmt.Errorf("error reading source: %v", err)
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Format:    v.opts.Encoder.String(),
		Source:    v.String(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (v *vault) String() string {
	return "vault"
}

func (v *vault) Watch() (source.Watcher, error) {
	w, err := newWatcher(v.prefix, v.String(), v.stripPrefix, v)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)

	// use default config
	config := api.DefaultConfig()

	// check if there are any addrs
	a, ok := options.Context.Value(addressKey{}).(string)
	if ok {
		if strings.HasPrefix(a, "https://") || strings.HasPrefix(a, "http://") {
			config.Address = a
		} else {
			addr, port, err := net.SplitHostPort(a)
			if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
				port = "8200"
				addr = a
				config.Address = fmt.Sprintf("https://%s:%s", addr, port)
			} else if err == nil {
				config.Address = fmt.Sprintf("https://%s:%s", addr, port)
			}
		}
	}

	// create the client
	client, _ := api.NewClient(config)

	prefix := DefaultPrefix
	sp := ""

	f, ok := options.Context.Value(prefixKey{}).(string)
	if ok {
		prefix = f
	}

	if b, ok := options.Context.Value(stripPrefixKey{}).(bool); ok && b {
		sp = prefix
	}

	return &vault{
		prefix:      prefix,
		stripPrefix: sp,
		addr:        config.Address,
		opts:        options,
		client:      client,
	}
}
