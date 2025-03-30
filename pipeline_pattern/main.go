package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// generate creates a channel that produces random integers between 0 and 999.
// It takes a context.Context as an argument to allow for cancellation of the generation process.
// The function runs a goroutine that continuously sends random integers to the channel until
// the context is done, at which point it closes the channel.
//
// Parameters:
//
//	ctx - a context.Context that can be used to signal cancellation of the generation process.
//
// Returns:
//
//	A read-only channel of integers that will produce random values until the context is cancelled.
func generate(ctx context.Context) <-chan int {
	out := make(chan int)
	flag := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				flag = 1
				break
			case out <- rand.Intn(1000):
			}
			if flag == 1 {
				close(out)
				break
			}
		}
	}()
	return out
}

// square takes a channel of integers as input and returns a channel of integers.
// It reads values from the input channel, squares each value, and sends the squared
// values to the output channel. The output channel is closed once all input values
// have been processed.
//
// Parameters:
// - channel: a read-only channel of integers from which values are received.
//
// Returns:
// - A read-only channel of integers containing the squared values.
func square(channel <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for val := range channel {
			out <- val * val
		}
		close(out)
	}()
	return out
}

// main is the entry point of the application. It creates a context with a timeout of 5 seconds,
// generates a sequence of numbers, squares them, and prints the results to the standard output.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	out := square(generate(ctx))

	for val := range out {
		fmt.Println(val)
	}
}
