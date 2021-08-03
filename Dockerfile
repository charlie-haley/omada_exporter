FROM golang:alpine3.14
RUN mkdir /app
WORKDIR /app

COPY . .
RUN go build -o main .

EXPOSE 9202
CMD ["/app/main"]
