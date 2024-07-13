# go-gprc-anomaly-detector

Client-server application that streams frequencies to the client that detects anomalies and puts them into PostgreSQL database.

## Launch

`docker-compose.yml` will allow you to build the client and the server as binaries to your machine, launch PostgreSQL database on port **5432** and launch server on **localhost:8888**. Just do:

```shell
docker compose up -d
```

To use client just launch binary from the **build** directory:

```shell
./client
```

## Options

### Server

- Command-line arguments
  - _optional_ address `string` - Server host and port, for example `localhost:8888`. Default: `0.0.0.0:8888`

### Client

- Command-line arguments
  - _optional_ `address` `string` - Server host and port, for example `localhost:8888`. Default: `0.0.0.0:8888`
  - _optional_ `k` `float` - Coeffitient for the detector to detect anomalies that only `k * STD` far from expected mean. Default: `1`
- Environment variables
  - _optional_ `POSTGRES_HOST` - PostgreSQL database host. Default: `localhost`
  - _optional_ `POSTGRES_USER` - PostgreSQL database user. Default: `postgres`
  - _optional_ `POSTGRES_PASSWORD` - PostgreSQL user password. Default: `postgres`
  - _optional_ `POSTGRES_DATABASE` - PostgreSQL database to connect to. Default: `postgres`
  - _optional_ `POSTGRES_PORT` - PostgreSQL port of the host. Default: `5432`

It is not necessary to create table for the anomalies in the database, client will be do it automatically.
