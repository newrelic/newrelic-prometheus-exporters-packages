#!/bin/bash

mkdir -p /etc/powerdns/pdns.d

PDNSVARS=`echo ${!PDNSCONF_*}`
touch /etc/powerdns/pdns.conf

PDNSCONF_GMYSQL_USER=pdns
PDNSCONF_GMYSQL_DBNAME=pdns
PDNSCONF_GMYSQL_PASSWORD=pdnspw


for var in $PDNSVARS; do
  varname=`echo ${var#"PDNSCONF_"} | awk '{print tolower($0)}' | sed 's/_/-/g'`
  value=`echo ${!var} | sed 's/^$\(.*\)/\1/'`
  echo "$varname=$value" >> /etc/powerdns/pdns.conf
done

if [ ! -z $PDNSCONF_API_KEY ]; then
  cat >/etc/powerdns/pdns.d/api.conf <<EOF
api=yes
webserver=yes
webserver-address=0.0.0.0
webserver-allow-from=0.0.0.0/0
EOF

fi

mysqlcheck() {
  # Wait for MySQL to be available...
  COUNTER=20
  until mysql -h mysql -u "$PDNSCONF_GMYSQL_USER" -p"$PDNSCONF_GMYSQL_PASSWORD" -e "show databases" 2>/dev/null; do
    echo "WARNING: MySQL still not up. Trying again..."
    sleep 10
    let COUNTER-=1
    if [ $COUNTER -lt 1 ]; then
      echo "ERROR: MySQL connection timed out. Aborting."
      exit 1
    fi
  done

  count=`mysql -h mysql -u "$PDNSCONF_GMYSQL_USER" -p"$PDNSCONF_GMYSQL_PASSWORD"  -e "select count(*) from information_schema.tables where table_type='BASE TABLE' and table_schema='dbns';" | tail -1`
  if [ "$count" == "0" ]; then
    echo "Database is empty. Importing PowerDNS schema..."
    mysql -h mysql -u "$PDNSCONF_GMYSQL_USER" -p"$PDNSCONF_GMYSQL_PASSWORD"  "$PDNSCONF_GMYSQL_DBNAME" < /usr/local/share/pdns/schema.mysql.sql && echo "Import done."
  else
    echo "database already exists..."
  fi
}

mysqlcheck

if [ "$SECALLZONES_CRONJOB" == "yes" ]; then
  cat > /etc/crontab <<EOF
PDNSCONF_API_KEY=$PDNSCONF_API_KEY
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
# m  h dom mon dow user    command
0,30 *  *   *   *  root    /usr/local/bin/secallzones.sh > /var/log/cron.log 2>&1
EOF
  ln -sf /proc/1/fd/1 /var/log/cron.log
  cron -f &
fi

# Start PowerDNS
# same as /etc/init.d/pdns monitor
echo "Starting PowerDNS..."

if [ "$#" -gt 0 ]; then
  exec /usr/sbin/pdns_server "$@"
else
  exec /usr/sbin/pdns_server --daemon=no --guardian=no --control-console --loglevel=9
fi
