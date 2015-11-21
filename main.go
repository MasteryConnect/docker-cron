package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
	"github.com/robfig/cron"
)

var running bool
var mu = &sync.Mutex{}

func main() {
	app := cli.NewApp()
	app.Name = "docker-cron"
	app.Usage = "used to run shell commands at specified intervals / times, based on cron syntax"
	app.Version = "1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Daniel Baldwin",
		},
	}
	app.Copyright = `
The MIT License (MIT)

Copyright (c) 2015 MasteryConnect

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
	`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "seconds",
			Value:  "*",
			Usage:  "seconds: 0-59, */10",
			EnvVar: "DC_SECS",
		},
		cli.StringFlag{
			Name:   "minutes",
			Value:  "*",
			Usage:  "minutes: 0-59, */10",
			EnvVar: "DC_MINS",
		},
		cli.StringFlag{
			Name:   "hours",
			Value:  "*",
			Usage:  "hours: 0-23, */10",
			EnvVar: "DC_HOURS",
		},
		cli.StringFlag{
			Name:   "day-of-month",
			Value:  "*",
			Usage:  "day of month: 1-31",
			EnvVar: "DC_DOM",
		},
		cli.StringFlag{
			Name:   "months",
			Value:  "*",
			Usage:  "month: 1-12 or JAN-DEC, */10",
			EnvVar: "DC_MONTHS",
		},
		cli.StringFlag{
			Name:   "day-of-week",
			Value:  "*",
			Usage:  "day of week: 0-6 or SUN-SAT",
			EnvVar: "DC_DOW",
		},
		cli.BoolFlag{
			Name:   "sync-jobs",
			Usage:  "should the jobs be run one at a time (true), or whenever they are scheduled",
			EnvVar: "DC_SYNC",
		},
	}
	app.Action = func(con *cli.Context) {
		// Vars
		var command string

		// Checks
		if len(con.Args()) > 0 {
			command = strings.Join(con.Args(), " ")
		} else {
			log.Fatal("Not enough args, need a command to run.")
			cli.ShowAppHelp(con)
		}

		// Ensure handling of SIGTERM and Interrupt
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		signal.Notify(signalChan, syscall.SIGTERM)
		go func() {
			<-signalChan
			os.Exit(1)
		}()

		// Setup cron job
		c := cron.New()
		schedule := strings.Join([]string{
			con.String("seconds"),
			con.String("minutes"),
			con.String("hours"),
			con.String("day-of-month"),
			con.String("months"),
			con.String("day-of-week"),
		}, " ")
		log.Printf("Setup cron to run on schedule: %s\n", schedule)
		c.AddFunc(schedule, func() {
			// Run one at a time, syncFlag=true, or whenever scheduled
			syncFlag := con.Bool("sync-jobs")
			if runJob(syncFlag) {
				defer jobDone(syncFlag)

				log.Printf("Running cron on schedule: %s\n", schedule)

				cmd := exec.Command("sh", "-c", command)

				setupStdout(cmd)
				setupStderr(cmd)

				err := cmd.Start()
				if err != nil {
					log.Fatal("Error running command", err)
				}

				err = cmd.Wait()
				if err != nil {
					log.Fatal("Error waiting for command", err)
				}
			} else {
				log.Println("A job is already running. The sync-jobs flag is true so we only run one at a time")
			}
		})
		c.Start()

		// Hold and let the cron job run
		for {
			time.Sleep(30 * time.Second)
		}
	}

	app.Run(os.Args)
}

func setupStdout(cmd *exec.Cmd) {
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdoutpipe for command", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()
}

func setupStderr(cmd *exec.Cmd) {
	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderrpipe for command", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Printf("ERR: %s\n", scanner.Text())
		}
	}()
}

// Should we run a job. If syncFlag is false, then we always run a job even if
// there is already one running. If syncFlag is true, then we only run a job
// if one is not already running
func runJob(syncFlag bool) bool {
	if syncFlag {
		mu.Lock()
		defer mu.Unlock()
		if running {
			return false
		} else {
			running = true
			return true
		}
	}
	// Always run, even if there is already one running
	return true
}

func jobDone(syncFlag bool) {
	// We only need to change the running state if we have syncrhonous job runs
	if syncFlag {
		mu.Lock()
		defer mu.Unlock()
		running = false
	}
}
