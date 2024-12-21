FROM golang:1.22.2 

EXPOSE 8080

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .


RUN go build -o main ./cmd/server

CMD ["./main"]