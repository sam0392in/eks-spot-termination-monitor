# -------- STAGE 1: Module download (cached) --------
FROM golang:1.24-alpine AS deps

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# -------- STAGE 2: Build --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN go build -o app ./cmd && \
    mv app /home

#---------STAGE 3-----------------

FROM golang:1.24-alpine

WORKDIR /home

ARG CONFIG_ENV=dev

COPY --from=builder /home/app ./app

EXPOSE 3000

CMD ["./app"]