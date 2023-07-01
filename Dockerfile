FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o receipt-processor-challenge ./cmd

CMD ["./receipt-processor-challenge"]