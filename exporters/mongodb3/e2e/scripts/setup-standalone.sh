#!/bin/bash 

mongodb1=`getent hosts ${MONGO1} | awk '{ print $1 }'`

port=${PORT:-27017}
DB="E2EStandaloneDB"
USER="root"
PASS="pass12345"

echo "Waiting for startup.."
until mongosh --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done

echo "Started.."

function noise_generator() {
      MAX_KEYS=10000
      function key() { KEY=$((RANDOM % $MAX_KEYS)); }
      while true; do
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.serverStatus())" > /dev/null
        key
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.E2EStandaloneFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.E2EStandaloneFirstCollection.deleteOne( { _id: $KEY } ))" > /dev/null
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.E2EStandaloneFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.E2EStandaloneSecondFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongosh $DB --host ${mongodb1}:${port} -u $USER -p $PASS --authenticationDatabase admin --eval "printjson(db.E2EStandaloneSecondFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        sleep $((RANDOM % 10))
      done
}

noise_generator
