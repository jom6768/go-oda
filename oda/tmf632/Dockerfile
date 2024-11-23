# Build Stage for tmf632
FROM golang:1.23-alpine AS builder-tmf632

WORKDIR /app

# Copy tmf632 source code and Go modules
COPY ./oda/tmf632 ./
COPY ./go.mod ./go.sum ./

# Install Go dependencies
RUN go mod download

# Build the Go binary for tmf632
RUN CGO_ENABLED=0 go build -o ./tmf632

# Final image for tmf632
FROM alpine:latest AS final-tmf632

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy the tmf632 binary from the builder-tmf632 stage
COPY --from=builder-tmf632 /app/tmf632 ./tmf632

EXPOSE 8632

CMD ["./tmf632"]