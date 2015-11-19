package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/codegangsta/cli"
	"github.com/robfig/cron"
)

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
			log.Printf("Running cron on schedule: %s\n", schedule)

			cmd := exec.Command("sh", "-c", command)

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

			err = cmd.Start()
			if err != nil {
				log.Fatal("Error running command", err)
			}

			err = cmd.Wait()
			if err != nil {
				log.Fatal("Error waiting for command", err)
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
