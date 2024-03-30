FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM scratch

WORKDIR /root

COPY --from=builder /main .

EXPOSE 8080

CMD ["/root/main"]
