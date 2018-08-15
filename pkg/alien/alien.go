package alien

import (
	"errors"
	"github.com/dangrier/alien/pkg/probe"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os/signal"
	"os"
)

type Alien struct {
	init       bool
	probes     map[*probe.Probe]bool
	results    chan probe.Result
	processing sync.Mutex

	stop chan time.Time
}

func New() *Alien {
	a := &Alien{
		init:       true,
		probes:     make(map[*probe.Probe]bool),
		results:    make(chan probe.Result),
		processing: sync.Mutex{},
		stop:       make(chan time.Time),
	}
	return a
}

func (a *Alien) AddProbe(p *probe.Probe) error {
	if !a.init {
		return errors.New("not initialised") // TODO: make const error
	}

	a.processing.Lock()
	defer a.processing.Unlock()

	a.probes[p] = true

	return p.Run()
}

// Run is the event loop, which blocks until Stop is called
func (a *Alien) Run() {
	// TODO: reconsider what event loop is needed for
	// Protect against uninitialised structs
	if !a.init {
		return
	}

	http.Handle("/metrics", promhttp.Handler())
	srv := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	go srv.ListenAndServe()

	for {
		select {
		case t := <-a.stop:
			// Stop requested

			logrus.Infof("Stop requested at: %v", t)

			srv.Close()

			a.processing.Lock()
			for p, _ := range a.probes {
				p.Stop()
			}
			a.processing.Unlock()
			return
		}
	}

	/*for res <-a.results {
		// Lock when actively processing to prevent
		// configuration reload breaking things
		a.processing.Lock()
		defer a.processing.Unlock() // TODO: don't defer in for loop

		if res.Error != nil {
			// There was an error
			// TODO: handle the error properly
			logrus.WithError(res.Error).Error("Probe error")
			continue
		}

		// TODO: add returned result to prometheus metric
	}
	*/
}

func (a *Alien) ListenForTermination() {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Kill, os.Interrupt)
		<-ch
		logrus.Info("Termination signal received") //TODO: remove
		a.stop <- time.Now()
	}()
}