#!/bin/bash

i=0
while true; do
    echo Message $i | ./opt/mqm/samp/bin/amqsputc Q1
    ((i++))
    sleep .1
done