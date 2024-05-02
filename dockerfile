FROM golang:latest

WORKDIR /app

ADD lamport .

RUN go build -o server .

CMD ["./server"]

EXPOSE 9000