package probe

import (
	"errors"
	"time"
)

type Option func(*Probe) error

func WithAuthBasic(user string, pass string) Option {
	return func(p *Probe) error {
		return errors.New("not implemented") // TODO: implement
	}
}

func WithAuthHeader(header string, value string) Option {
	return func(p *Probe) error {
		return errors.New("not implemented") // TODO: implement
	}
}

func WithFrequency(frequency time.Duration) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.Stop()
		p.ticker = time.NewTicker(frequency)
		return nil
	}
}

func WithPayload(payload string) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.payload = payload
		return nil
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(p *Probe) error {
		p.processing.Lock()
		defer p.processing.Unlock()
		p.client.Timeout = timeout
		return nil
	}
}

func WithSuccessFilter(filter ResultFilter) Option {
	return func(p *Probe) error {
		if p.success != nil {
			return ErrFilterAlreadySet
		}
		p.success = filter
		return nil
	}
}
