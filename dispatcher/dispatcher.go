package dispatcher

import (
	"fmt"
	"io"
	"strconv"
	"sync"
)

// Message serves for sending data to registered channel
type Message struct {
	Reciever uint64
	Data     []byte
}

// Dispatcher manages chanels for communication
type Dispatcher struct {
	lock     sync.RWMutex
	items    map[uint64]chan []byte
	income   chan Message
	signal   chan bool
	errors   chan error
	isClosed bool
}

// New creates and initialises Dispatcher
func New() *Dispatcher {
	d := Dispatcher{
		lock:     sync.RWMutex{},
		items:    make(map[uint64]chan []byte),
		income:   make(chan Message),
		errors:   make(chan error),
		signal:   make(chan bool),
		isClosed: false,
	}
	go d.run()
	return &d
}

func (d *Dispatcher) run() {
	select {
	case x := <-d.income:
		fmt.Printf("In Run %v: %v\n", x.Reciever, x.Data)
		d.lock.RLock()
		fmt.Println("In run.lock " + strconv.FormatUint(x.Reciever, 10))
		ch, err := d.get(x.Reciever)
		if err != nil {
			d.errors <- err
		} else {
			ch <- x.Data
		}
		d.lock.RUnlock()
		fmt.Println("Run.Unlock " + strconv.FormatUint(x.Reciever, 10))
	case <-d.signal:
		return
	}
}

func (d *Dispatcher) hasItem(id uint64) bool {
	_, ok := d.items[id]
	return ok
}

func (d *Dispatcher) delete(id uint64) {
	if v, ok := d.items[id]; ok {
		close(v)
		delete(d.items, id)
	}
}

func (d *Dispatcher) add(id uint64) (<-chan []byte, error) {
	ch := make(chan []byte)
	d.items[id] = ch
	return ch, nil
}

// Register adds new channel to dispatcher and returns channel for writing
func (d *Dispatcher) Register(id uint64) (<-chan []byte, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isClosed {
		return nil, io.EOF
	}

	d.delete(id)
	return d.add(id)
}

func (d *Dispatcher) get(id uint64) (chan []byte, error) {
	if v, ok := d.items[id]; ok {
		return v, nil
	}
	return nil, io.ErrUnexpectedEOF
}

// Size return count of registered items
func (d *Dispatcher) Size() int {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if d.isClosed {
		return 0
	}
	return len(d.items)
}

// Delete close read channel and removes ID from register
func (d *Dispatcher) Delete(id uint64) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.isClosed {
		return
	}
	d.delete(id)
}

// Errors returns channel to get errors
func (d *Dispatcher) Errors() <-chan error {
	return d.errors
}

// Close closes all channels and free all resources
func (d *Dispatcher) Close() {
	d.lock.Lock()
	defer d.lock.Unlock()

	close(d.signal)
	for k := range d.items {
		d.delete(k)
	}
	close(d.errors)
	close(d.income)
	d.isClosed = true
}

// Send routes data to id's chanell
func (d *Dispatcher) Send(id uint64, data []byte) error {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.isClosed {
		return io.EOF
	}

	d.income <- Message{id, data}
	fmt.Printf("Sended for %v: %v\n", id, data)
	return nil
}
