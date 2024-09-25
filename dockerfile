FROM golang:1.22-alpine

WORKDIR /app

COPY . .

# Copy Go module files to utilize cache and avoid re-downloading if no changes are detected
COPY go.* ./
RUN go mod download
RUN go mod verify

RUN go build -o /tickets-api ./main.go

EXPOSE 8080

CMD ["/tickets-api"]
