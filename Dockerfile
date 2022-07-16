FROM golang:1.18-alpine

WORKDIR /usr/src

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
WORKDIR /usr/src/server
RUN go build -o /app

CMD ["/app"]
