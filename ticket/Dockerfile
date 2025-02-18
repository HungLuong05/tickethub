FROM golang:alpine

WORKDIR /ticket

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ticket .

EXPOSE 8002

CMD ["./ticket"]