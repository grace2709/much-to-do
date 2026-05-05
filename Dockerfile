# ---------- BUILD STAGE ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o muchtodo

# ---------- RUN STAGE ----------
FROM alpine:latest

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/muchtodo .

EXPOSE 8080

HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./muchtodo"]
