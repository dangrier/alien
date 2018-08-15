package probe

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"sync"
	"time"
	"github.com/sirupsen/logrus"
	"fmt"
	"bytes"
)

type Probe struct {
	init       bool
	processing sync.Mutex
	running    bool

	client *http.Client
	prom   *prometheus.CounterVec

	endpoint string
	method   string
	payload  string
	freq     time.Duration
	ticker   *time.Ticker
	success  ResultFilter

	stop chan time.Time
}

func New(endpoint string, options ...Option) (*Probe, error) {
	// Generate default struct values
	p := &Probe{
		init:     true,
		client:   http.DefaultClient,
		endpoint: endpoint,
		method:   "GET",
		payload:  "",
		ticker:   time.NewTicker(10 * time.Second),
		freq:     1 * time.Minute,
		stop:     make(chan time.Time),
	}

	for _, o := range options {
		if err := o(p); err != nil {
			return nil, err
		}
	}

	p.prom = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "alien_probe_count",
		Help: "Count of probes by endpoint and success",
	}, []string{
		"endpoint",
		"success",
	})
	prometheus.MustRegister(p.prom)

	return p, nil
}

func (p *Probe) Run() error {
	logrus.Info("Running")
	if !p.init {
		return ErrNotInitialised
	}

	if p.running {
		return ErrNotStopped
	}

	// Run the main loop in a new goroutine
	if err := p.Validate(); err != nil{
		return err
	}
	go p.run()

	return nil
}

// run is the main loop for a Probe, and is intended to
// be run concurrently in a separate goroutine.
func (p *Probe) run() {
	// Protect against uninitialised structs
	if !p.init {
		return
	}

	for {
		logrus.Info("Event loop") // TODO: remove
		select {
		case <-p.stop:
			// Cancellation has been requested
			p.ticker.Stop()
			p.running = false
			return

		case <-p.ticker.C:
			// TODO: finish/fix probe check

			p.processing.Lock()

			body:=bytes.NewBufferString(p.payload)
			req, err := http.NewRequest(p.method, p.endpoint, body)
			if err != nil {
				logrus.WithError(err).Error("Creating request")
				p.processing.Unlock()
				continue
			}
			res, err := p.client.Do(req)
			if err != nil {
				logrus.WithError(err).Error("Getting response")
				p.processing.Unlock()
				continue //TODO: add prometheus counter for errors
			}

			probeResult := &Result{
				Timestamp: time.Now(),
				Probe: p,
				Code: res.StatusCode,
				Headers: res.Header,
			}

			logrus.Info("Probe check simulation") // TODO: remove
			p.prom.WithLabelValues(p.endpoint, fmt.Sprintf("%t", p.success.Check(probeResult))).Inc() // TODO: update to use real values
			p.processing.Unlock()
		}
	}
}

// Stop sends a signal to stop further processing
// (after current processing if any), and stops the
// timer from running further.
//
// Stop **will block** until it is able to stop successfully.
func (p *Probe) Stop() error {
	logrus.Info("Stopping") // TODO: remove
	if !p.init {
		return ErrNotInitialised
	}

	if !p.running {
		return ErrNotRunning
	}

	p.stop <- time.Now()
	return nil
}

// Validate checks whether there are enough valid data to
// carry out a probe check. Returns nil if no problems, otherwise
// returns an Error with the reason for failure.
func (p *Probe) Validate() error {
	if !p.init {
		return ErrNotInitialised
	}

	if p.endpoint == "" {
		return ErrInvalidEndpoint
	}

	if p.method == "" {
		return ErrInvalidMethod
	}

	if p.freq <= 0 {
		return ErrInvalidFrequencyZero
	}

	if p.success == nil {
		return ErrInvalidSuccessFilterEmpty
	}

	return nil
}
