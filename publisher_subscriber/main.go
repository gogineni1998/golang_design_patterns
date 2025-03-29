package main

import (
	"fmt"
	"strconv"
)

type Subscriber struct {
	subscriberID int
	events       chan string
}

// publish sends a message to all subscribers.
// It takes a message of type string and a pointer to a slice of Subscriber.
// Each subscriber's events channel is used to deliver the message.
func publish(message string, subscribers *[]Subscriber) {
	for _, subscriber := range *subscribers {
		subscriber.events <- message
	}
}

// subscriberReceiver listens for events on the subscriber's event channel
// and prints each received event along with the subscriber's ID.
// It runs in a separate goroutine for each subscriber to handle events concurrently.
func subscriberReceiver(subscriber *Subscriber) {
	for event := range subscriber.events {
		fmt.Println("Subscriber "+strconv.Itoa(subscriber.subscriberID)+" received:", event)
	}
}

// main is the entry point of the application. It provides a simple command-line interface
// for publishing messages and subscribing to them. The user can choose to publish a message,
// subscribe to receive messages, or exit the application. When a message is published, it is
// sent to all active subscribers. When a subscriber is added, a new goroutine is started to
// handle incoming messages for that subscriber.
func main() {
	var option int
	var message string
	subscribers := []Subscriber{}
	for {
		fmt.Println("Enter 1 : Publish \n 2 : Subscribe \n Enter any thing except 1 or 2 to stop")
		fmt.Scanf("%d\n", &option)
		if option == 1 {
			fmt.Println("Enter Message")
			fmt.Scanln(&message)
			publish(message, &subscribers)
		} else if option == 2 {
			subscribers = append(subscribers, Subscriber{
				subscriberID: len(subscribers) + 1,
				events:       make(chan string),
			})
			go subscriberReceiver(&subscribers[len(subscribers)-1])
		} else {
			for _, subscriber := range subscribers {
				close(subscriber.events)
			}
			break
		}
	}
}
