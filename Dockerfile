FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum file
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod Download

# Copy the source code
COPY . .

# Build the main server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o flight-api ./cmd/server/main.go

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate/main.go

# Create a minimal production image
FROM alpine:3.22

WORKDIR /app

# Install required packages
RUN apk --no-cache add tz-data

# Copy binaries from the builder stage
COPY --from=builder /app/flight-api /app/flight-api
COPY --from=builder /app/migrate /app/migrate

# Copy migrations
COPY --from=builder /app/db/migrations /app/db/migrations

# Add a non-root user and switch to it
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the entry point to the main server by default
ENTRYPOINT [ "/app/flight-api" ]