package vault

import (
	"context"
	"time"

	"github.com/micro/go-config/source"
)

type addressKey struct{}
type pollKey struct{}
type prefixKey struct{}
type stripPrefixKey struct{}

// WithAddress sets the consul address
func WithAddress(a string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, addressKey{}, a)
	}
}

// WithPrefix sets the key prefix to use
func WithPrefix(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, prefixKey{}, p)
	}
}

// PollInterval specifies the interval on which to poll vault since it does not support watches
func PollInterval(t time.Duration) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, pollKey{}, t)
	}
}

// StripPrefix indicates whether to remove the prefix from config entries, or leave it in place.
func StripPrefix(strip bool) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, stripPrefixKey{}, strip)
	}
}
