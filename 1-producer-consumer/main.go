//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweets chan *Tweet, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(tweets)
	}()
	for {
		tweet, err := stream.Next()
		if errors.Is(err, ErrEOF) {
			break
		}
		tweets <- tweet
	}
}

func consumer(tweets chan *Tweet, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	tweets := make(chan *Tweet, 100)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Producer
	go producer(stream, tweets, wg)
	wg.Add(1)

	// Consumer
	go consumer(tweets, wg)
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
