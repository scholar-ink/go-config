package vault

import (
	"errors"
	"time"

	"github.com/micro/go-config/source"
)

type watcher struct {
	v           *vault
	cs          *source.ChangeSet
	poll        time.Duration
	name        string
	stripPrefix string

	ch   chan *source.ChangeSet
	exit chan bool
}

func newWatcher(key, name, stripPrefix string, v *vault) (source.Watcher, error) {
	cs, err := v.Read()
	if err != nil {
		return nil, err
	}

	pollInterval := time.Minute

	p, ok := v.opts.Context.Value(pollKey{}).(time.Duration)
	if ok {
		pollInterval = p
	}

	w := &watcher{
		cs:          cs,
		v:           v,
		poll:        pollInterval,
		name:        name,
		stripPrefix: stripPrefix,
		ch:          make(chan *source.ChangeSet),
		exit:        make(chan bool),
	}

	go w.watch()

	return w, nil
}

func (w *watcher) watch() {
	t := time.NewTicker(w.poll)
	defer t.Stop()

	for {
		select {
		case <-w.exit:
			return
		case <-t.C:
			cs, err := w.v.Read()
			if err != nil {
				continue
			}
			if cs.Checksum == w.cs.Checksum {
				continue
			}
			w.cs = cs
			w.ch <- cs
		}
	}
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case cs := <-w.ch:
		return cs, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return nil
	default:
		close(w.exit)
	}
	return nil
}
