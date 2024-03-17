FROM golang:alpine3.19

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

RUN go build -o main .
CMD ["./main"]
