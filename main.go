package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func client(message chan string, ch chan time.Time, interval time.Duration) {
	for {
		message <- "Current time: "
		ch <- time.Now()
		time.Sleep(time.Second * interval)
	}


}

func server(message chan string, ch chan time.Time)  {
	for {
		select {
		case msg, ok := <- message:
			if !ok {
				return
			}
			fmt.Println(msg)

		case timestamp, ok := <- ch:
			if !ok {
				return
			}
			fmt.Println(timestamp)
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
		fmt.Println(err.Error())
	}
	seconds := bytes[0]
	interval := time.Duration(seconds)

	wg := sync.WaitGroup{}
	defer wg.Done()
	wg.Add(2)

	ch := make(chan time.Time)
	message := make(chan string)

	go client(message, ch, interval)
	go server(message, ch)

	wg.Wait()
}
