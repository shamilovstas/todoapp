FROM golang:1.23.4 AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todoapp-backend cmd/main.go

FROM alpine:latest
EXPOSE 8080
COPY --from=build /app/todoapp-backend .
CMD ["./todoapp-backend"]