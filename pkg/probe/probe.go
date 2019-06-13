package probe

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Probe is an abstraction for a set of basic procedures
// to carry out to check an endpoint, and to manage the
// conditions which indicate a successful probe.
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

	failureActions []Action
	successActions []Action

	logger logrus.StdLogger

	stop chan time.Time
}

// New is a Probe constructor which makes sure defaults are applied
// then allows for variadic functional options to be provided
// to further configure the probe.
func New(endpoint string, options ...Option) (*Probe, error) {
	// Generate default struct values
	p := &Probe{
		init:     true,
		client:   http.DefaultClient,
		endpoint: endpoint,
		method:   "GET",
		payload:  "",
		ticker:   time.NewTicker(10 * time.Second),
		freq:     10 * time.Second,
		stop:     make(chan time.Time),
	}

	p.logger = log.New(os.Stdout, fmt.Sprintf("%s: ", p), 0)

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

	p.logger.Printf("%s: Registered prometheus metric collector", p)

	return p, nil
}

// Run actually makes the Probe function by starting the
// ticker. Will return an ErrNotStopped error if already running.
//
// Run also triggers run() - the main loop for a Probe's functionality,
// and starts the first Trigger() call to execute a probe.
func (p *Probe) Run() error {
	if !p.init {
		return ErrNotInitialised
	}

	if p.running {
		return ErrNotStopped
	}

	// Run the main loop in a new goroutine
	if err := p.Validate(); err != nil {
		return err
	}

	go p.run()
	p.running = true

	p.Trigger()

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
		select {
		case <-p.stop:
			// Cancellation has been requested
			p.ticker.Stop()
			p.running = false
			p.logger.Printf("%s: Stopped", p)
			return

		case <-p.ticker.C:
			p.Trigger()
		}
	}
}

// SetLogger sets the logger for the probe
func (p *Probe) SetLogger(l logrus.StdLogger) {
	p.logger = l
	l.Printf("%s: Set logger", p)
}

// Stop sends a signal to stop further processing
// (after current processing if any), and stops the
// timer from running further.
//
// Stop **will block** until it is able to stop successfully.
func (p *Probe) Stop() error {
	if !p.init {
		return ErrNotInitialised
	}

	p.logger.Printf("%s: Stopping", p)

	if !p.running {
		p.logger.Printf("%s: Already stopped", p)
		return ErrNotRunning
	}

	p.stop <- time.Now()
	return nil
}

// String implements the Stringer interface
func (p *Probe) String() string {
	return fmt.Sprintf("Probe<%s '%s' every %s>", p.method, p.endpoint, p.freq)
}

// Trigger a probe to do a check now
func (p *Probe) Trigger() error {
	p.processing.Lock()
	defer p.processing.Unlock()

	if !p.init {
		return ErrNotInitialised
	}

	if err := p.Validate(); err != nil {
		return err
	}

	p.logger.Printf("%s: Triggered...", p)

	body := bytes.NewBufferString(p.payload)
	req, err := http.NewRequest(p.method, p.endpoint, body)
	if err != nil {
		p.logger.Printf("%s: failed: %v", p, err)
		p.prom.WithLabelValues(p.endpoint, "false").Inc() // Record error as success=false
		return err
	}
	res, err := p.client.Do(req)
	if err != nil {
		p.logger.Printf("%s: failed: %v", p, err)
		p.prom.WithLabelValues(p.endpoint, "false").Inc() // Record error as success=false
		return err
	}

	probeResult := &Result{
		Timestamp: time.Now(),
		Probe:     p,
		Code:      res.StatusCode,
		Headers:   res.Header,
	}

	p.logger.Printf("%s: Completed", p)

	if p.success.Check(probeResult) {
		p.prom.WithLabelValues(p.endpoint, "true").Inc()
		for _, a := range p.successActions {
			a(*probeResult)
		}
	} else {
		p.prom.WithLabelValues(p.endpoint, "false").Inc()
		for _, a := range p.failureActions {
			a(*probeResult)
		}
	}

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
