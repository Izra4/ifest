FROM golang:1.23
WORKDIR /app
COPY . .
COPY go.mod go.sum ./
RUN go mod tidy
RUN go build -o main cmd/main.go
EXPOSE 3000
CMD ["./main"]