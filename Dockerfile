# ========================
# Single Stage Build & Run
# ========================
FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o server .

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