FROM debian
COPY ./node-allocation-exporter /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/node-allocation-exporter"]
