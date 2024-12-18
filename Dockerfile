FROM golang:1.22.2 

EXPOSE 8080

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
COPY ./config/appsettings.json /app/config/appsettings.json

RUN go build -o main ./cmd/server
RUN go test ./content

CMD ["./main"]