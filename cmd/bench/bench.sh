#!/bin/sh

config=$1
data=$2
queries=$3

go build -tags annoy -o bench cmd/bench/main.go

for algo in brute_force brute_force_blas annoy sanny
do
  if [ $algo = "sanny" ]; then
    for inner_algo in annoy ngt
    do
      ./bench -data $data -algo $algo -inner-algo $inner_algo -config $config -test-size $queries
    done
  else
      ./bench -data $data -algo $algo -config $config -test-size $queries
  fi
done
