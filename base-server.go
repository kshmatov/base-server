package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/kshmatov/base-server/dispatcher"
)

func runner(id int, d *dispatcher.Dispatcher, wg *sync.WaitGroup) {
	defer wg.Done()

	ch, err := d.Register(uint64(id))
	if err != nil {
		return
	}
	for x := range ch {
		fmt.Printf("%v: %v", id, string(x))
	}
}

func main() {
	ids := make([]int, 10)
	wg := new(sync.WaitGroup)
	dsp := dispatcher.New()

	for i := 0; i < len(ids); i++ {
		wg.Add(1)
		go runner(i, dsp, wg)
	}

	go func(dsp *dispatcher.Dispatcher) {
		for x := range dsp.Errors() {
			fmt.Printf("Error: %v\n", x)
		}
	}(dsp)

	rand.Seed(89)
	for i := 0; i < 100; i++ {
		e := dsp.Send(uint64(rand.Intn(11)), []byte("Message-"+strconv.Itoa(i)))
		if e != nil {
			fmt.Printf("Send: %v\n", e)
		}
	}

	dsp.Close()
	fmt.Println("Finished")
}
