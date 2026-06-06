FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o omada-exporter .

FROM alpine:3.18
COPY --from=builder /app/omada-exporter /usr/bin/omada-exporter
CMD ["/usr/bin/omada-exporter"]
