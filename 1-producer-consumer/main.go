//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, ch chan<- *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return
		}
		ch <- tweet
	}
}

func consumer(ch <-chan *Tweet) {
	for t := range ch {
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

	// Producer
	var tweetchan = make(chan *Tweet)
	go func() { producer(stream, tweetchan) }()

	// Consumer
	consumer(tweetchan)

	fmt.Printf("Process took %s\n", time.Since(start))
}
