FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .

#required by go-sqlite3.
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o otp-server .

FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN useradd -r -u 1001 -g root otpserver
WORKDIR /app
COPY --from=builder /build/otp-server .
RUN chown otpserver /app/otp-server
RUN mkdir -p /data && chown otpserver /data
USER otpserver
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/app/otp-server"]
