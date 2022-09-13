#!/bin/bash

# Forcing the client authentication to be required for DEV.APP.SVRCONN
cat << EOF > init.file
ALTER CHANNEL(DEV.APP.SVRCONN) CHLTYPE(SVRCONN) SSLCAUTH(REQUIRED)
ALTER CHANNEL(DEV.ADMIN.SVRCONN) CHLTYPE(SVRCONN) SSLCAUTH(REQUIRED)
DISPLAY CHANNEL(DEV.ADMIN.SVRCONN)
DISPLAY CHANNEL(DEV.APP.SVRCONN)
EOF

# running the server
bash -c "sleep 10 ; runmqsc -f init.file" & runmqdevserver
