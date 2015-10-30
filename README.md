# docker-cron
A simple one liner to run a command on a cron schedule.

## Overview
`docker-cron` was born out of wanting a simple one liner way to keep a docker container up with a command being run at specific intervals. Integrating cron into a container in a generic way always presented itself as a much longer and harder process than was really worth it to just run a simple script. So we built `docker-cron` to make getting a simple single purpose docker container up a running faster.

## Code Example

Print out "hello world" every 10 seconds
```
docker-cron --seconds=*/10 echo hello world
```

Run file every 15 minutes
```
docker-cron --seconds=0 --minutes=*/15 ./test.sh
```

To see the help text
```
docker-cron --help
```
