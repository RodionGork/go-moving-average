# Go Moving Average

A kind of exercise: program which processes multiple requests which either sent new values for some metrics
or request to report moving average for the latest period of time of fixed length.

Start it:

    make build

    build/server

Test manually:

    echo -n "power 200" | nc -u -w1 127.0.0.1 1090
    echo -n "voltage 350" | nc -w1 127.0.0.1 1080
    echo -n "power 250" | nc -w1 127.0.0.1 1080
    echo -n "power" | nc -u -w1 127.0.0.1 1090

First three invocations simply send forth the values for given metrics and won't
print anything. The last queries average power during last minute (it should print `225`).
