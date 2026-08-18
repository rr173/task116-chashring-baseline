FROM golang:1.26.3

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

RUN go build -mod=vendor ./...

CMD ["bash"]
