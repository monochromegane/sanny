#!/bin/sh

config=$1
data=$2
queries=$3

go build -tags annoy -o bench cmd/bench/main.go

for algo in brute_force brute_force_blas annoy sanny
do
  ./bench -data $data -algo $algo -config $config -test-size $queries
done
