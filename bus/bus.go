package bus

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
)

type (
	Pump interface {
		// Open a new stream of data from the given ID
		Open(id uint64) error

		// Close indicates that a given ID is no longer interested in receiving messages
		// it might be called more than once for the same id
		Close(id uint64) error

		// PumpIn a new batch of messages from a given ID
		Push(id uint64, payload ...Payload) error

		// PumpOut messages to the given ID (including broadcast ones, ie.: 0),
		// if an ID is passed to close, PumpOut should return an error
		Pop(id uint64) ([]Payload, error)
	}

	B struct {
		outputLock sync.RWMutex
		output     map[uint64]chan Payload
		workers    chan signal

		handlerID   uint64
		handlerLock sync.RWMutex
		handlers    map[uint64]handler
	}

	Handler interface {
		Handle(output Transmitter, from uint64, p Payload) error
	}

	HandlerFunc func(output Transmitter, from uint64, p Payload) error

	Payload struct {
		Value json.RawMessage `json:"value"`
		Meta  struct {
			Kind string     `json:"kind"`
			Path []NodePath `json:"path"`
		} `json:"meta"`
	}

	NodePath struct {
		ID      string `json:"id"`
		Trigger string `json:"etID"`
		Tag     string `json:"tag"`
	}

	Transmitter interface {
		Transmit(to uint64, p Payload)
	}

	handler func(done chan signal, out Transmitter, from uint64, p Payload)

	signal struct{}
)

func New() *B {
	return &B{
		output:   make(map[uint64]chan Payload),
		workers:  make(chan signal, 1000),
		handlers: make(map[uint64]handler),
	}
}

func (hf HandlerFunc) Handle(output Transmitter, from uint64, p Payload) error {
	return hf(output, from, p)
}

func (b *B) Push(from uint64, payload ...Payload) error {
	b.handlerLock.RLock()
	defer b.handlerLock.RUnlock()
	for _, msg := range payload {
		for _, h := range b.handlers {
			b.workers <- signal{}
			go h(b.workers, b, from, msg)
		}
	}
	return nil
}

func (b *B) RegisterHandlerFunc(fn func(output Transmitter, from uint64, p Payload) error) uint64 {
	return b.RegisterHandler(HandlerFunc(fn))
}

func (b *B) RegisterHandler(fn Handler) uint64 {
	nid := atomic.AddUint64(&b.handlerID, 1)
	b.handlerLock.Lock()
	b.handlers[nid] = func(done chan signal, t Transmitter, from uint64, p Payload) {
		defer func() { <-done }()
		fn.Handle(t, from, p)
	}
	b.handlerLock.Unlock()
	return nid
}

func (b *B) Close(id uint64) error {
	b.outputLock.Lock()
	ch, found := b.output[id]
	if found {
		delete(b.output, id)
	}
	b.outputLock.Unlock()
	if found {
		close(ch)
		return nil
	}
	return errors.New("already closed")
}

func (b *B) Open(id uint64) error {
	b.outputLock.Lock()
	_, found := b.output[id]
	if found {
		b.outputLock.Unlock()
		return errors.New("duplicate")
	}
	ch := make(chan Payload, 100)
	b.output[id] = ch
	b.outputLock.Unlock()
	return nil
}

func (b *B) Transmit(to uint64, p Payload) {
	b.outputLock.RLock()
	if to == 0 {
		// broadcast
		for _, o := range b.output {
			select {
			case o <- p:
				return
			default:
				// discard, slow consumer
			}
		}
	} else {
		ch := b.output[to]
		if ch == nil {
			return
		}
		select {
		case ch <- p:
		default:
			// discard, slow consumer
		}
	}
	b.outputLock.RUnlock()
}

func (b *B) Pop(to uint64) ([]Payload, error) {
	b.outputLock.RLock()
	ch, found := b.output[to]
	b.outputLock.RUnlock()
	if !found {
		return nil, errors.New("not found")
	}
	p, open := <-ch
	if !open {
		return nil, errors.New("closed")
	}
	acc := []Payload{p}
	for {
		select {
		case p, open := <-ch:
			if open {
				acc = append(acc, p)
			}
		default:
			return acc, nil
		}
	}
}
