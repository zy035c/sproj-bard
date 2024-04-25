FROM golang:latest

WORKDIR /lamport

COPY . .

RUN go build -o server .

CMD ["./server", "-port=9090", "-ps=8000,9000,10000"]