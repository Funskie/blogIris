FROM golang:alpine as builder

LABEL maintainer="Funskie <tusty9292@gmail.com>"

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Strat new stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. And copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
