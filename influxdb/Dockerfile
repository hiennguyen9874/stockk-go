ARG INFLUXDB_VERSION=2.6-alpine
FROM influxdb:${INFLUXDB_VERSION}

COPY ./createSecondBucket.sh /docker-entrypoint-initdb.d
RUN chmod +x /docker-entrypoint-initdb.d/createSecondBucket.sh
