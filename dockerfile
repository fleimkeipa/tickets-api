FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /tickets-api ./main.go

EXPOSE 8080

CMD ["/tickets-api"]