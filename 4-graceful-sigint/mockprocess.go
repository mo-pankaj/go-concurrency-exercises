//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
// Your task is to edit `main.go`
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

// MockProcess for example
type MockProcess struct {
	isRunning bool
}

// Run will start the process
func (m *MockProcess) Run(ctx context.Context) {
	m.isRunning = true

	fmt.Print("Process running..")
	for {
		select {
		case <-ctx.Done():
			m.Stop()
		default:
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}
	}
}

// Stop tries to gracefully stop the process, in this mock example
// this will not succeed
func (m *MockProcess) Stop() {
	if !m.isRunning {
		log.Fatal("Cannot stop a process which is not running")
	}

	fmt.Print("\nStopping process..")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Print("\nKilling process..")
		os.Exit(0)
	}()
	fmt.Print(".")
	time.Sleep(10 * time.Second)
	os.Exit(0)
}
