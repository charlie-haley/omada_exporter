FROM golang:alpine3.14

COPY omada-exporter /

CMD ["/omada-exporter"]
