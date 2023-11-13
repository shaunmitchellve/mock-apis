FROM golang:1.21.1-alpine3.18
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/function ./
FROM alpine:3.18
COPY --from=0 /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]