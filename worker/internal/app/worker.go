package app

import (
	"math/rand/v2"
	"time"

	"github.com/paulja/go-work/worker/internal/ports"
)

var _ ports.Worker = (*Worker)(nil)

type Worker struct {
	Cancel chan interface{}
}

func (w *Worker) Start(id, payload string) error {
	// this is where the work actually get's done. The payload
	// could be a command or some actual work. But this where
	// if the worker had a consumer library where could offload
	// the actual work outside of the `go-work` machinery.
	//
	// In this mock-up a random timer between 30-60 seconds is
	// created to pretend to be doing some work.

	min, max := uint(30), uint(60)
	num := rand.UintN(max-min) + min

	cancel := make(chan interface{})
	select {
	case <-time.After(time.Duration(num) * time.Second):
		close(cancel)
		cancel = nil
	case <-cancel:
		cancel = nil
	}
	w.Cancel = cancel

	return nil
}

func (w *Worker) Stop(id string) error {
	if w.Cancel != nil {
		close(w.Cancel)
	}
	return nil
}
