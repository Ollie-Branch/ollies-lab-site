# Dockerfile
FROM golang:latest AS base

# Set working directory
WORKDIR /app

COPY . .

RUN go mod download && go build -o ollies-lab

EXPOSE 6767

CMD ["/app/ollies-lab"]

