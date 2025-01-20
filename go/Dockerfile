# Base Go environment
# -------------------
FROM golang:1.23-alpine as base
WORKDIR /hatchet

COPY go.mod go.sum ./

RUN go mod download

COPY /main.go ./main.go

# Go build environment
# --------------------
FROM base AS build-go

RUN go build -a -o ./bin/hatchet-worker .

# Deployment environment
# ----------------------
FROM alpine AS deployment

WORKDIR /hatchet

# openssl and bash needed for admin build
RUN apk update && apk add --no-cache gcc musl-dev openssl bash ca-certificates

COPY --from=build-go /hatchet/bin/hatchet-worker /hatchet/

CMD /hatchet/hatchet-worker
