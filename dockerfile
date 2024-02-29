FROM golang:1.21.0-alpine3.17
WORKDIR /app
COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY main.go ./main.go
COPY internal ./internal
RUN go build -o main .
CMD ["./main","-c","./config.yml"]