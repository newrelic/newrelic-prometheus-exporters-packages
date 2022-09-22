## E2E

Thanks to the `MQ_DEV` environment variable two channels are created, one for administration, the other for normal messaging:
 - `DEV.ADMIN.SVRCONN` - configured to only allow the admin user to connect into it. A user and password must be supplied.
 - `DEV.APP.SVRCONN` - does not allow administrative users to connect. Password is optional unless you choose a password for app users.

Consumer and producers connects to `DEV.APP.SVRCONN` using the `app` user. And the exporters uses the `admin` one to connect to `DEV.ADMIN.SVRCONN`.

In the mTLS scenario they also both uses the same certificate in the keystore to authenticate, but different `ccdt.json`. 

## TLS
One of the tests covers a mTLS scenario setting up server, consumer, producers and an exporters all leveraging mTLS and username and password to authenticate.

### keystore

The keystore is used by clients to store the certificate and connect to the server leveraging mTLS.

The pem certificates is signed by the ca of the server and stored under the trusted folder that is mounted in `/etc/mqm/pki/trust/0`
All pem certificates need to be added as well to the keystore mounted by the client and pointed by the environment variable `MQSSLKEYR=/key`

The client needs as well to indicate which certificate should be used when connecting. This is performed thanks to the `ccdt.json` file, with the `certificateLabel` property.

The following commands where used in order to create the keystore.
```bash
# To create the server.p12 that can be loaded into the keystore
$ openssl pkcs12 -export -out server.p12 -inkey my_private_key.pem  -in ./my_signed_cert.pem  -certfile ./exporters/ibmmq/e2e/tls/ca.crt

# create key store
$ runmqakm -keydb -create -db new.kdb -pw test -type pkcs12 -expire 0 -stash
# import certificate
$ runmqakm -cert -import -file server.p12 -pw test -type pkcs12 -target new.kdb -target_pw test -target_type cms
# relabel the certificate
$ runmqckm -cert -rename -db new.kdb -type pkcs12 -label cn=*,o=newrelic -new_label admin
# list the certificate store
$ runmqakm -cert -list -db key.kdb
# Certificates found
# * default, - personal, ! trusted, # secret key
# !       "l=san fransisco,c=us,cn=demo.mlopshub.com"
# -       admin
```

The keystore `key.kdb` is the only file needed to be mounted by clients and exporter together with the stashed password `key.sth` and the `ccdt.json` file that specify how the connection needs to be opened.

### ccdt.json

This is a json specifying how a connection needs to be opened. The `cipherSpecification` is used for TLS and `certificateLabel`
for mTLS to specify which certificate should be selected from the keystore and sent.

Notice that the `queueManager` specified should match the one in the amqsputc commands, for example `amqsputc`

### Server

The server has mTLS configured thanks to the option `SSLCAUTH` enabled at startup time in the `entrypoint-server-tls`

> If you want SSLCAUTH=REQUIRED This means that the client application must present its own personal 
> certificate to the queue manager for validation,then you must ensure you have a personal 
> certificate and its signers (if using CA-signed certs) in the applications keystore.

The certificate of the server is stored into `/etc/mqm/pki/keys/mykey` and the client one in `/etc/mqm/pki/trust/0`

### Simple TLS
You can also test the "simple" TLS to do so it is enough to remove from the server entrypoint:

```bash
ALTER CHANNEL(DEV.APP.SVRCONN) CHLTYPE(SVRCONN) SSLCAUTH(REQUIRED)
ALTER CHANNEL(DEV.ADMIN.SVRCONN) CHLTYPE(SVRCONN) SSLCAUTH(REQUIRED)
```

and from the `ccdt.json` file:
```bash
   "certificateLabel": "admin"
```
