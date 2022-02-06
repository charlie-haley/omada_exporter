FROM golang:alpine3.14

COPY omada_exporter /

CMD ["omada_exporter"]
