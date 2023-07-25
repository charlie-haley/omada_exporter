FROM golang:alpine3.14 AS build

WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY main.go ./
COPY cmd ./cmd
COPY pkg ./pkg
RUN go build -o /omada-exporter

FROM golang:alpine3.14
COPY --from=build /omada-exporter /omada-exporter
CMD "/omada-exporter"
