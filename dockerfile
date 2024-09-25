FROM golang:1.22-alpine

WORKDIR /app

COPY . .

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

RUN go build -o /tickets-api ./main.go

EXPOSE 8080

CMD ["/tickets-api"]