//
// timecat
//
// (c) 2016-2017 Chris J Arges <christopherarges@gmail.com>
//

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

func setup() {
	// Don't show signal names so we can display time properly.
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func cleanup(timestamp bool) {
	if !timestamp {
		// Add a newline.
		fmt.Println()
	}
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func timecat(args []string, timestamp bool) {

	// Setup tty
	setup()

	// On ctrl-c exit gracefully
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cleanup(timestamp)
		os.Exit(0)
	}()

	// Get current time.
	start := time.Now()

	// Roughly update at 30fps
	if !timestamp {
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
	}

	// Execute command if supplied, otherwise loop.
	if len(args) > 1 {
		cmd := exec.Command(args[0], args[1:]...)
		stdout, _ := cmd.StdoutPipe()
		cmd.Stderr = os.Stderr
		cmd.Start()

		// Buffer and print output including timestamps if requested
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			if timestamp {
				log.Printf(in.Text())
			} else {
				fmt.Printf(in.Text())
			}
		}
		if err := in.Err(); err != nil {
			log.Printf("error: %s", err)
		}

		// Show final time
		cmd.Wait()
		if !timestamp {
			s := fmt.Sprintf("%v", time.Since(start))
			fmt.Printf(s)
		}

		// Manually clean up
		cleanup(timestamp)
	} else {
		for {
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "timecat"
	app.Usage = "a better time utility"
	app.Version = "0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "timestamp",
		},
	}

	app.Action = func(c *cli.Context) error {
		timecat(c.Args(), c.Bool("timestamp"))
		return nil
	}

	app.Run(os.Args)
}

// vim: tabstop=4
