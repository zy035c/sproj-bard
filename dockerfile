FROM golang:latest

WORKDIR /lamport

COPY . .

RUN go build -o server .

CMD ["./server"]