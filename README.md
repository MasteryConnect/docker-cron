# docker-cron
A simple one liner to run a command on a cron schedule.

## Overview
`docker-cron` was born out of wanting a simple one liner way to keep a docker container up with a command being run at specific intervals. Integrating cron into a container in a generic way always presented itself as a much longer and harder process than was really worth it to just run a simple script. So we built `docker-cron` to make getting a simple single purpose docker container up and running faster.

## Code Example

Print out "hello world" every 10 seconds
```
docker-cron --seconds=*/10 echo hello world
```

Run file every 15 minutes
```
docker-cron --seconds=0 --minutes=*/15 ./test.sh
```

By default a new process is run every scheduled interval, even if the previously run process is still running. If you don't want that i.e. if you want only one process to run at any given time, use the --sync-jobs (env DC_SYNC=true) option. If the previously scheduled process is still running when the next interval time elapses, a new process won't be created.

Example. Run every 5 seconds. However with the command sleeping (running) for 6 seconds, every second scheduled interval will not run the command.
```
docker-cron --seconds=*/5 --sync-jobs sleep 6
```

To see the help text
```
docker-cron --help

NAME:
   docker-cron - used to run shell commands at specified intervals / times, based on cron syntax

USAGE:
   docker-cron [global options] command [command options] [arguments...]
   
VERSION:
   1.0
   
AUTHOR(S):
   Daniel Baldwin 
   
COMMANDS:
   help, h      Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --seconds "*"        seconds: 0-59, */10 [$DC_SECS]
   --minutes "*"        minutes: 0-59, */10 [$DC_MINS]
   --hours "*"          hours: 0-23, */10 [$DC_HOURS]
   --day-of-month "*"   day of month: 1-31 [$DC_DOM]
   --months "*"         month: 1-12 or JAN-DEC, */10 [$DC_MONTHS]
   --day-of-week "*"    day of week: 0-6 or SUN-SAT [$DC_DOW]
   --sync-jobs          should the jobs be run one at a time (true), or whenever they are scheduled [$DC_SYNC]
   --help, -h           show help
   --version, -v        print the version
```
