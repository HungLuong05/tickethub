FROM golang:alpine

WORKDIR /event

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o event .

EXPOSE 8001
# EXPOSE 50051

CMD ["./event"]