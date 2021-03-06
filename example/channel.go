package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

func main() {
	// set up channel
	ch := make(chan interface{})

	// create channel-based flow and set an OnValue trigger.
	// The flow will send notifications every time the defined value 3
	// is send through the channel.
	flow := flow.New(&filters.OnValue{3}, &flow.Chan{ch})
	defer flow.Close()

	// syncrhoniztion
	var wg sync.WaitGroup

	// publisher: random numbers to be added in irregular intervals
	wg.Add(2)

	go publisher(1, ch, &wg)
	go publisher(2, ch, &wg)

	// subscribers
	wg.Add(2)

	go subscriber(1, flow, &wg)
	go subscriber(2, flow, &wg)

	wg.Wait()
}

func publisher(id int, ch chan interface{}, wg *sync.WaitGroup) {
	var counter int
	for {
		val := rand.Intn(4)
		if val == 3 {
			counter++
		}
		ch <- val
		sleep := rand.Intn(2)

		fmt.Printf("Publisher %d sends: value = %v, counts = %d, sleeps = %d sec \n", id, val, counter, sleep)
		time.Sleep(time.Duration(sleep) * time.Second)

		if counter >= 5 {
			break
		}
	}
	wg.Done()
}

func subscriber(id int, monitor observer.Observer, wg *sync.WaitGroup) {
	sub := monitor.Subscribe()
	for i := 1; i < 10; i++ {
		<-sub.C()
		fmt.Printf("Subscriber %d got notified: value = %v, counts = %d\n", id, sub.Value(), i)
	}
	wg.Done()
}
