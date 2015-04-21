#!/bin/bash

# use ./startServers <num of servers to start>

numServers=$1
i=0
port=10000
while [ $i -lt $numServers ]
do
	./server -clientAddress 127.0.0.1:$((port+i)) -logfile node -ID $i & 
	i=$((i+1))
done


