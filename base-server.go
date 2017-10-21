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
		fmt.Printf("Ups %v: %v\n", id, err)
		return
	}
	fmt.Printf("Registered %v\n", id)
	for x := range ch {
		fmt.Printf("Recieved for %v: %v\n", id, string(x))
	}
}

func main() {
	ids := make([]int, 2)
	wg := new(sync.WaitGroup)
	dsp := dispatcher.New()

	for i := 0; i < len(ids); i++ {
		wg.Add(1)
		go runner(i, dsp, wg)
	}

	go func(dsp <-chan error) {
		for x := range dsp {
			fmt.Printf("Error: %v\n", x)
		}
	}(dsp.Errors())

	rand.Seed(89)
	for i := 0; i < 10; i++ {
		fmt.Printf("Iteration %v\n", i)
		e := dsp.Send(uint64(rand.Intn(11)), []byte("Message-"+strconv.Itoa(i)))
		if e != nil {
			fmt.Printf("Send error: %v\n", e)
		}
	}

	dsp.Close()
	wg.Wait()
	fmt.Println("Finished")
}
