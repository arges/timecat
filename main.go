//
// timecat
//
// (c) 2016-2017 Chris J Arges <christopherarges@gmail.com>
//

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func setup() {
	// Don't show signal names so we can display time properly.
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func cleanup() {
	fmt.Println()
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func main() {
	// Get command line arguments
	args := os.Args[1:]

	// Setup tty
	setup()

	// On ctrl-c exit gracefully
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cleanup()
		os.Exit(0)
	}()

	// Get current time.
	start := time.Now()

	// Roughly update at 30fps
	go func() {
		c := time.Tick(32 * time.Millisecond)
		for _ = range c {
			s := fmt.Sprintf("%v", time.Since(start))
			fmt.Printf(s)
			for i := 0; i < len(s); i++ {
				fmt.Printf("\b")
			}
		}
	}()

	// Execute command if supplied, otherwise loop.
	if len(args) > 1 {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		// Show final time
		s := fmt.Sprintf("%v", time.Since(start))
		fmt.Printf(s)

		// Manually clean up
		cleanup()
	} else {
		for {
		}
	}
}

// vim: tabstop=4
