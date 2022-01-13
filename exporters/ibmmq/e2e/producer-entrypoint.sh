#!/bin/bash

i=0
while true; do
    echo Message $i | ./opt/mqm/samp/bin/amqsputc DEV.QUEUE.1
    ((i++))
    sleep .1
done