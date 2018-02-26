package netvm

import (
	"net"
	"net/http"
	"sync"
)

const (
	// HydrationIDHeader is the header used to pass the hydration ID to the deadhead
	HydrationIDHeader = "Hydration-ID"
)

// OnDeadheadError called when an error occurs inside a deadhead
type OnDeadheadError func(err error)

// HydratedFunc is an http handler that can also signal a headead of an error
type HydratedFunc func(w http.ResponseWriter, r *http.Request) error

// Hydrater is anything that can fill a func with an ID
type Hydrater interface {
	Hydrate(id string) (HydratedFunc, error)
}

// ServeDeadhead serves a deadhead request
func ServeDeadhead(l net.Listener, hydrater Hydrater, onError OnDeadheadError) error {
	deadhead := &Deadhead{
		hydrater: hydrater,
		onError:  onError,
	}

	return http.Serve(l, deadhead)
}

// Deadhead represents an empty husk of a function, waiting to breath life into itself
type Deadhead struct {
	hydrater    Hydrater
	hydrateOnce sync.Once
	f           HydratedFunc
	onError     OnDeadheadError
}

func (d *Deadhead) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	d.hydrateOnce.Do(func() {
		id := r.Header.Get(HydrationIDHeader)
		if len(id) == 0 {
			return
		}

		d.f, err = d.hydrater.Hydrate(id)
		if err != nil {
			d.onError(err)
			return
		}
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if d.f == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.Header.Del(HydrationIDHeader)

	err = d.f(w, r)
	if err != nil {
		d.onError(err)
	}
}
