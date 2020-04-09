# Loki Query Load Generator

This project is a simple golang application to generate query load to a Loki server.

To build the app run `make build`.

`./loki-gen-load --url=http://localhost:3100/` to start querying a Loki server at `http://localhost:3100/`

You can use the `--verbose` option to print the rate per second of request sent.

To build a docker image use `make build-image`.

To run the image use: `docker run ctovena/loki-gen-load:0.1 --url=http://localhost:3100`
