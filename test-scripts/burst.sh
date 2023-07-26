#!/bin/bash

for i in {1..3}
do
  time METRICS=20000 BATCHES=10 BATCH_SIZE=100000 go run client.go &
done