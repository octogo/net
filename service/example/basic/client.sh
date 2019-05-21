#!/bin/bash

i=0

while (( i <= 100 )); do
    telnet localhost 31337 &>/dev/null &
    i=$(( i + 1 ))
done