package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// Broker is a helper for sending and receiving messages.
type Broker struct {
	client *redis.Client
}

var ctx = context.Background()

// New is a constructor for Broker.
func New() Broker {
	return Broker{
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}


func (b Broker) receiveTask(key string) string {
	result, _ := b.client.RPop(key).Result()
	fmt.Println(result)
	return result
}

func handleTask(value string) {
	fmt.Println(value)
	fmt.Println("Sum of ", value, " is: ", value)
}
func main() {
	TaskKey := "email:new"
	broker := New()
	for {
		time.Sleep(time.Second)
		taskValue := broker.receiveTask(TaskKey)
		if len(taskValue) > 0 {
			go handleTask(taskValue)
		}
	}
}