# Build Stage for tmf629
FROM golang:1.23-alpine AS builder-tmf629

WORKDIR /app

# Copy tmf629 source code and Go modules
COPY ./oda/tmf629 ./
COPY ./go.mod ./go.sum ./

# Install Go dependencies
RUN go mod download

# Build the Go binary for tmf629
RUN CGO_ENABLED=0 go build -o ./tmf629

# Final image for tmf629
FROM alpine:latest AS final-tmf629

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy the tmf629 binary from the builder-tmf629 stage
COPY --from=builder-tmf629 /app/tmf629 ./tmf629

EXPOSE 8629

CMD ["./tmf629"]
