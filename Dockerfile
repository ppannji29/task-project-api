# ========================
# Build Stage
# ========================
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o server .

# ========================
# Run Stage
# ========================
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 3000
CMD ["./server"]


# # ========================
# # Build Stage
# # ========================
# FROM golang:1.23 AS builder

# WORKDIR /app

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . ./

# RUN go build -o server .

# EXPOSE 3000

# CMD ["./server"]