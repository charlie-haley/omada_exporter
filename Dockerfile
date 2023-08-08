FROM golang:alpine3.14

COPY omada-exporter /usr/bin/omada-exporter

CMD ["/usr/bin/omada-exporter"]
