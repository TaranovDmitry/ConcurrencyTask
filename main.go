package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type message struct {
	time time.Time
	msg  string
}


func client(chMessage chan message, interval time.Duration, group *sync.WaitGroup, ctx context.Context) {
	defer group.Done()
	for {
		time.Sleep(interval)
		select {
			case chMessage <- message{
				time: time.Now(),
				msg:  "Current time:",
		}:
		case <- ctx.Done():
			fmt.Println("Time out")
			return
		}

	}
}

func server(chMessage chan message, group *sync.WaitGroup, ctx context.Context) {
	defer group.Done()
	for {
		select {
		case message, ok := <-chMessage:
			if !ok {
				return
			}
			fmt.Println(message.msg, message.time)
		case <- ctx.Done():
			return

		default:
			time.Sleep(time.Second * 1)
		}
	}
}

func main() {
	file, err := os.Open("config")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(string(bytes))
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	messages := make(chan message, 5)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	
	go client(messages, duration, &wg, ctx)
	go server(messages, &wg, ctx)
	wg.Wait()
}
