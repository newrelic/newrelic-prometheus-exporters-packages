#!/bin/bash

while true; do
    echo ${MQ_APP_PASSWORD} | /opt/mqm/samp/bin/amqsgetc DEV.QUEUE.1 QM1
done
