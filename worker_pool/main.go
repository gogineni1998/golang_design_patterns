package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// generate reads integers from a file named "task.txt" and sends them to a channel.
// It returns a receive-only channel of integers. The function opens the file, reads
// each line, converts it to an integer, and sends it to the channel until the end of
// the file is reached or an error occurs. The channel is closed once all integers have
// been sent.
//
// Returns:
//
//	<-chan int: A receive-only channel that will yield integers read from the file.
func generate() <-chan int {
	out := make(chan int)
	file, err := os.Open("task.txt")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	go func() {
		for {
			task, err := reader.ReadString('\n')
			if err == io.EOF {
				fmt.Println(err)
				break
			}
			task = strings.TrimSuffix(task, "\n")
			task = strings.TrimSuffix(task, "\r")
			task_needed, err := strconv.Atoi(task)
			if err != nil {
				fmt.Println(err)
				break
			}
			out <- task_needed
		}
		close(out)
	}()
	return out
}

// taskWorkers takes a channel of tasks and a number of workers,
// and returns a channel that will receive the results of the tasks.
// Each task is processed concurrently by a limited number of workers.
//
// Parameters:
// - tasks: a read-only channel of integers representing the tasks to be processed.
// - workers: an integer specifying the maximum number of concurrent workers.
//
// Returns:
// - A read-only channel of integers containing the results of the processed tasks.
func taskWorkers(tasks <-chan int, workers int) <-chan int {
	output := make(chan int)
	worker_pool := make(chan bool, workers)
	wg := sync.WaitGroup{}
	go func() {
		for task := range tasks {
			worker_pool <- true
			wg.Add(1)
			go func() {
				output <- task * task
				<-worker_pool
				wg.Done()
			}()
		}
		wg.Wait()
		close(worker_pool)
		close(output)
	}()
	return output
}

// main is the entry point of the application. It prompts the user to enter the number of workers,
// generates a set of tasks, and processes those tasks using the specified number of workers.
// It also measures and prints the time taken to complete the task processing.
func main() {
	var workers string
	fmt.Println("Worker Poll Pattern :-")
	fmt.Println("Enter the number of workers :-")
	fmt.Scanf("%s\n", &workers)
	timer := time.Now()
	number_of_workers, err := strconv.Atoi(workers)
	if err != nil {
		panic(err)
	}
	tasks := generate()
	processes_output := taskWorkers(tasks, number_of_workers)
	for po := range processes_output {
		fmt.Println(po)
	}
	fmt.Println("Time taken :- ")
	fmt.Println(time.Since(timer))
}
