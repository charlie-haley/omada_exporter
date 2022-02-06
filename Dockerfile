FROM golang:alpine3.14

COPY omada_exporter /usr/bin/omada_exporter

CMD ["/usr/bin/omada_exporter"]
