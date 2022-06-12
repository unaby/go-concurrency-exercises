//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) AddTime(seconds int64) int64 {
	return atomic.AddInt64(&u.TimeUsed, int64(seconds))
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	timer := time.NewTicker(time.Second)
	done := make(chan interface{}, 1)
	// Short circuit if they're already out of time
	if atomic.LoadInt64(&u.TimeUsed) >= 10 {
		return false
	}
	go processRequest(done, process)
	for {
		select {
		case <-timer.C:
			if i := u.AddTime(1); i >= 10 {
				if !u.IsPremium {
					return false
				}
			}
		case <-done:
			return true
		}
	}
	return true
}

func processRequest(done chan interface{}, process func()) {
	process()
	done <- true
}

func main() {
	RunMockServer()
}
