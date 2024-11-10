# Build Stage for tmf632
FROM golang:1.23.1 AS builder-tmf632

WORKDIR /app

# Copy tmf632 source code and Go modules
COPY ./oda/tmf632 ./oda/tmf632
COPY ./go.mod ./go.sum ./

# Install Go dependencies
RUN go mod download
RUN go mod tidy

# Build the Go binary for tmf632
RUN go build -o /app/tmf632 ./oda/tmf632

# Build Stage for tmf669
FROM golang:1.23.1 AS builder-tmf669

WORKDIR /app

# Copy tmf669 source code and Go modules
COPY ./oda/tmf669 ./oda/tmf669
COPY ./go.mod ./go.sum ./

# Install Go dependencies
RUN go mod download
RUN go mod tidy

# Build the Go binary for tmf669
RUN go build -o /app/tmf669 ./oda/tmf669

# Final image for tmf632
FROM alpine:latest AS final-tmf632

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy the tmf632 binary from the builder-tmf632 stage
COPY --from=builder-tmf632 /app/tmf632 /app/tmf632

EXPOSE 8081

CMD ["./tmf632"]

# Final image for tmf669
FROM alpine:latest AS final-tmf669

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy the tmf669 binary from the builder-tmf669 stage
COPY --from=builder-tmf669 /app/tmf669 /app/tmf669

EXPOSE 8082

CMD ["./tmf669"]