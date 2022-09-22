#!/bin/bash

i=0
while true; do
    echo "${MQ_APP_PASSWORD}
      Message $i
      " | /opt/mqm/samp/bin/amqsputc DEV.QUEUE.1 QM1

    ((i++))
    sleep 10
done
