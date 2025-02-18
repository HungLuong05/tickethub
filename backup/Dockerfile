FROM golang:alpine

WORKDIR /backup

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o backup .

# EXPOSE 8001

CMD ["./backup"]