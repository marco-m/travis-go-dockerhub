FROM golang:1.13-alpine as builder

ENV GO111MODULE=on

RUN apk --no-cache add curl
RUN curl --silent --location https://taskfile.dev/install.sh | sh

WORKDIR /code

#
# Optimize downloading of dependencies only when they are needed.
# This requires to _first_ copy only these two files, run `go mod download`,
# and _then_ copy the rest of the source code.
#
COPY go.mod go.sum ./
RUN go mod download

#
# Build.
#

COPY . .

RUN task build

#
# The final image
#

FROM alpine

# You need this if your executable uses the Go HTTP client code.
#RUN apk --no-cache add ca-certificates

COPY --from=builder /code/bin/* /bin/
