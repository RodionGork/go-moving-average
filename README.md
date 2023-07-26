# Go Moving Average

A kind of exercise: program which processes multiple requests which either sent new values for some metrics
or request to report moving average for the latest period of time of fixed length.

**Start it:**

    make build

    build/server

setting `DEBUG_COUNTER` environment variable (to anything) will cause server to produce log messages about
what is going on, which may be useful for manual testing, but in speed run may affect performance.

**Test manually:**

Server takes commands via plain `TCP` or `UDP` connection, as simple text lines in form

    metric_name value

Multiple commands could be sent at once, with lines separated by newline character.

If the `value` is omitted, then command is treated as request to report current average (during last minute).

Easiest way to try it is by using `telnet`:

    telnet 127.0.0.1 1080
    power 1000
    voltage 300
    power 2000
    voltage
    300 1
    power
    1500 2
    ... wait about 1 minute ...
    power
    0 0

Here values `300`, `1500` and `0` are sent as response by server. They are followed by the total number of
values participating in calculation of given average value. E.g. `1500` is average over `2` values (we remember
them to be `1000` and `2000`).

To send commands from command line it is very convenient to use `netcat`:

    echo "power 200" | nc -u -w1 127.0.0.1 1090
    echo "voltage 350" | nc -w1 127.0.0.1 1080
    echo "power 250" | nc -w1 127.0.0.1 1080
    echo "power" | nc -u -w1 127.0.0.1 1090

First three invocations simply send forth the values for given metrics and won't
print anything. The last queries average power during last minute (it should print `225`).

**Test automatically**

There is a small script to send multiple commands, it could be invoked like this:

    go run test-scripts/client.go

By default it send `1` batch of `1000` commands using `10` distinct names for metrics. These values could be
changed by using environment variables:

    ADDRESS=127.0.0.1:1080 BATCHES=100 BATCH_SIZE=10000 METRICS=30000 go run test-scripts/sclient.go

after loading is over, it is easy to check for results using `netcat` again, metric names are `m***` where
stars represent integer value, `0`-based up to the total number of metrics, i.e.

    echo "m0\nm29000" | nc -u -w1 127.0.0.1 1090

Here is also `test-scripts/burst.sh` - small bash script allowing to run several processes with this tool.
