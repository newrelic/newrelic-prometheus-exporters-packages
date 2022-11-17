#!/bin/bash

mongodb1=`getent hosts ${MONGOS} | awk '{ print $1 }'`
mongodb11=`getent hosts ${MONGO11} | awk '{ print $1 }'`
mongodb12=`getent hosts ${MONGO12} | awk '{ print $1 }'`
mongodb21=`getent hosts ${MONGO21} | awk '{ print $1 }'`
mongodb22=`getent hosts ${MONGO22} | awk '{ print $1 }'`

port=${PORT:-27017}
DB="E2EShardDB"

echo "Waiting for startup.."
until mongo --host ${mongodb1}:${port} --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done

echo "Started.."

mongo --host ${mongodb1}:${port} <<EOF
   sh.addShard( "${RS1}/${mongodb11}:${PORT1},${mongodb12}:${PORT2}" );
   sh.addShard( "${RS2}/${mongodb21}:${PORT1},${mongodb22}:${PORT2}" );
   sh.status();
EOF

function noise_generator() {
      MAX_KEYS=10000
      function key() { KEY=$((RANDOM % $MAX_KEYS)); }
      while true; do
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.serverStatus())" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EShardFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EShardFirstCollection.deleteOne( { _id: $KEY } ))" > /dev/null
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EShardFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EShardSecondCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EShardSecondCollection.insertOne( { _id: $KEY } ))" > /dev/null
        sleep $((RANDOM % 10))
      done
}

noise_generator
