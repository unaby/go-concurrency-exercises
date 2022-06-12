//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
)

func main() {
	// Create a process

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	shutdown := make(chan bool)
	proc := MockProcess{}

	go func() {
		<-sigChan
		go proc.Stop()
		<-sigChan
		shutdown <- true
	}()
	// Run the process (blocking)

	go func() {
		proc.Run()
		shutdown <- true
	}()

	<-shutdown

}
