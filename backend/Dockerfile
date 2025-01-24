FROM golang:1.23.4
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todoapp-backend cmd/main.go
EXPOSE 8080
CMD ["./todoapp-backend"]