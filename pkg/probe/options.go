package probe

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// Option provides a way of configuring a probe using
// variadic parameters when calling probe.New()
//
// The Option is applied after the initial creation of
// the struct instance, and can only be set at init time.
type Option func(*Probe) error

// WithAuthBasic sets basic authentication credentials
// for the probe
func WithAuthBasic(user string, pass string) Option {
	return func(p *Probe) error {
		return errors.New("not implemented") // TODO: implement
	}
}

// WithHeader sets the given header to the given value for
// the probe
func WithHeader(header string, value string) Option {
	return func(p *Probe) error {
		return errors.New("not implemented") // TODO: implement
	}
}

// WithFrequency sets the rate at which probes will be conducted
//
// If not used, the default is 10 seconds
func WithFrequency(frequency time.Duration) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.Stop()
		p.ticker = time.NewTicker(frequency)
		p.freq = frequency
		return nil
	}
}

// WithLogger sets the probe's logger to use
func WithLogger(l logrus.StdLogger) Option {
	return func(p *Probe) error {
		p.logger = l
		return nil
	}
}

// WithMethod sets the probe's HTTP method
func WithMethod(method string) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.method = method
		return nil
	}
}

// WithPayload sets a body to send in a request
func WithPayload(payload string) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.payload = payload
		return nil
	}
}

// WithClient sets the outgoing request HTTP client
//
// If not used, the http.DefaultClient is used
func WithClient(timeout time.Duration) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.client.Timeout = timeout
		return nil
	}
}

// WithSuccessFilter sets the conditions considered to
// be a successful probe, given as a composed ResultFilter.
//
// See the `ResultFilter` for more detail, but ResultFilter
// is a set of conditions about the response code or content.
func WithSuccessFilter(filter ResultFilter) Option {
	return func(p *Probe) error {
		if p.success != nil {
			return ErrFilterAlreadySet
		}
		p.success = filter
		return nil
	}
}

// Action represents a callback function which does something
// with a Result
type Action func(Result)

// OnFailure adds the given Action func to the list of actions
// to perform when a probe failure (non-success) is identified
func OnFailure(a Action) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.failureActions = append(p.failureActions, a)
		return nil
	}
}

// OnSuccess adds the given Action func to the list of actions
// to perform when a probe success is identified
func OnSuccess(a Action) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.successActions = append(p.successActions, a)
		return nil
	}
}
