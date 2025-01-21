FROM golang:alpine

WORKDIR /auth

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth .

EXPOSE 8000
EXPOSE 50051

CMD ["./auth"]