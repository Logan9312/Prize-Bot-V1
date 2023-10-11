FROM golang:1.21.3-alpine AS build_base

RUN apk add --no-cache git
RUN apk add build-base

# Set the Current Working Directory inside the container
WORKDIR /tmp/app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./main .

FROM alpine:3.15
RUN apk add ca-certificates

COPY --from=build_base /tmp/app/main /main

# Run the startup script
CMD ["/main"]
