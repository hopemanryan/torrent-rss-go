
FROM golang:1.19-alpine3.16 
WORKDIR /app

RUN apk update
RUN apk add  gcc python3-dev g++
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN cd ./pgk && go build -v -o ./

# CMD [ "go run ./cmd/main.go" ]
CMD ["./pkg/pkg"]
