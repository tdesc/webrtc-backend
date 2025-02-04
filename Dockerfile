# Build stage
FROM golang:1.18 as builder

WORKDIR /app
COPY go.mod .
# If you have a go.sum file, copy it as well.
# COPY go.sum .
RUN go mod download

COPY . .

# Build the server binary.
RUN go build -o server .

# Final stage
FROM golang:1.18
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["./server"]
