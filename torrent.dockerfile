
FROM golang:1.19-alpine3.16 as base 
WORKDIR /app

RUN apk update
RUN apk add  gcc python3-dev g++
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN cd ./pkg && go build -v -o ./torrent

CMD ["./pkg/torrent"]
