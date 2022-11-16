#!/bin/bash 

mongodb1=`getent hosts ${MONGO1} | awk '{ print $1 }'`
mongodb2=`getent hosts ${MONGO2} | awk '{ print $1 }'`
arbiter=`getent hosts ${ARBITER} | awk '{ print $1 }'`

port=${PORT:-27017}
DB="E2EReplicasetDB"

echo "Waiting for startup.."
until mongo --host ${mongodb1}:${port} --eval 'quit(db.runCommand({ ping: 1 }).ok ? 0 : 2)' &>/dev/null; do
  printf '.'
  sleep 1
done

echo "Started.."

function setup_servers() {
    echo "setup servers"
    mongo --host ${mongodb1}:${port} <<EOF
    var cfg = {
        "_id": "${RS}",
        "protocolVersion": 1,
        "members": [
            {
                "_id": 0,
                "host": "${mongodb1}:${port}"
            },
            {
                "_id": 1,
                "host": "${mongodb2}:${port}"
            }
        ]
    };
    rs.initiate(cfg, { force: true });
    rs.reconfig(cfg, { force: true });

    rs.addArb("${arbiter}:${port}")
EOF
}

function noise_generator() {
      MAX_KEYS=10000
      function key() { KEY=$((RANDOM % $MAX_KEYS)); }
      while true; do
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.serverStatus())" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EReplicasetFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EReplicasetFirstCollection.deleteOne( { _id: $KEY } ))" > /dev/null
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EReplicasetFirstCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EReplicasetSecondCollection.insertOne( { _id: $KEY } ))" > /dev/null
        key
        mongo $DB --host ${mongodb1}:${port} --eval "printjson(db.E2EReplicasetSecondCollection.insertOne( { _id: $KEY } ))" > /dev/null
        sleep $((RANDOM % 10))
      done
}

setup_servers
noise_generator
