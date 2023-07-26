package main
// run it with `go run client.go`

import (
    "math/rand"
    "net"
    "os"
    "strconv"
)

func main() {
    address := "127.0.0.1:1080"
    if a := os.Getenv("ADDRESS"); a != "" {
        address = a
    }
    metrics := 10
    if m := os.Getenv("METRICS"); m != "" {
        metrics, _ = strconv.Atoi(m)
    }
    batches := 1
    if b := os.Getenv("BATCHES"); b != "" {
        batches, _ = strconv.Atoi(b)
    }
    batchSize := 1000
    if s := os.Getenv("BATCH_SIZE"); s != "" {
        batchSize, _ = strconv.Atoi(s)
    }
    
    for j := 0; j < batches; j++ {
        conn, err := net.Dial("tcp", address)
        if err != nil {
            println(err.Error())
            os.Exit(1)
        }
        for i := 0; i < batchSize; i++ {
            v := 200 + (i%metrics) + rand.Int() % 100
            conn.Write([]byte("m" + strconv.Itoa(i%metrics) + " " + strconv.Itoa(v) + "\n"))
        }
        conn.Close()
    }
}
