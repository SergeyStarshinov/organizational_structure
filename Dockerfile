FROM golang:alpine

WORKDIR /app

COPY ./go.mod ./

RUN go mod download

COPY . . 

EXPOSE 8080

RUN go build ./cmd/main.go

CMD ["./main"]