package alien

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dangrier/alien/pkg/probe"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Alien is the controller for a set of configured probes
type Alien struct {
	init       bool
	probes     map[*probe.Probe]bool
	results    chan probe.Result
	processing sync.Mutex

	metricsAddress  string
	metricsPort     int
	metricsEndpoint string

	logger logrus.StdLogger

	stop chan time.Time
}

// New is a constructor for an Alien and handles
// the setup of maps/channels/mutex.
func New() *Alien {
	a := &Alien{
		init:            true,
		probes:          make(map[*probe.Probe]bool),
		results:         make(chan probe.Result),
		processing:      sync.Mutex{},
		metricsAddress:  "",
		metricsEndpoint: "/metrics",
		metricsPort:     8080,
		logger:          log.New(os.Stdout, "Alien: ", 0),
		stop:            make(chan time.Time),
	}
	return a
}

// AddProbe tells an Alien to manage the provided Probe
func (a *Alien) AddProbe(p *probe.Probe) error {
	if !a.init {
		return ErrNotInitialised
	}

	a.logger.Printf("Adding probe %s", p)

	a.processing.Lock()
	defer a.processing.Unlock()

	a.probes[p] = true

	return p.Run()
}

// RemoveProbe tells an Alien to stop managing the provided Probe
func (a *Alien) RemoveProbe(p *probe.Probe) error {
	if !a.init {
		return ErrNotInitialised
	}

	a.logger.Printf("Removing probe %s", p)

	a.processing.Lock()
	defer a.processing.Unlock()

	if _, ok := a.probes[p]; !ok {
		return ErrProbeNotFound
	}

	p.SetLogger(log.New(os.Stdout, "", 0))

	delete(a.probes, p)
	return nil
}

// Run is the event loop, which intentionally blocks until Stop is called
func (a *Alien) Run() {
	// Protect against uninitialised structs
	if !a.init {
		return
	}

	a.listenForTermination()

	a.logger.Printf("Starting metrics handler %s:%d%s...", a.metricsAddress, a.metricsPort, a.metricsEndpoint)
	http.Handle(a.metricsEndpoint, promhttp.Handler())
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.metricsAddress, a.metricsPort),
		Handler: nil,
	}
	go srv.ListenAndServe()

	for {
		select {
		case <-a.stop:
			// Stop requested
			a.logger.Println("Stop signal received, terminating...")

			srv.Close()

			a.processing.Lock()
			for p := range a.probes {
				p.Stop()
			}
			a.processing.Unlock()
			return
		}
	}
}

// SetLogger sets the logger to use for all its probe logs
func (a *Alien) SetLogger(l logrus.StdLogger) {
	a.logger = l
	l.Println("Set logger for alien")
}

// listenForTermination opens a new goroutine which
// waits for an os.Kill or os.Interrupt and when received
// sends a value on the stop channel, which gracefully exits
func (a *Alien) listenForTermination() {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Kill, os.Interrupt)
		<-ch
		a.stop <- time.Now()
	}()
}
