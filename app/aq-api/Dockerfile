# Build Environment
FROM golang:1.21.3-alpine3.18 as build
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/
COPY go.mod go.sum main.go ./
# Setup other modules
RUN mkdir routes
COPY routes ./routes
RUN mkdir data
COPY data ./data
RUN go mod download
RUN go build -o /usr/local/bin/air-quality-mock-endpoint ./

# Production Environment
FROM alpine:3.18
COPY --from=build /usr/local/bin/air-quality-mock-endpoint /usr/local/bin
RUN mkdir config
COPY config /config
RUN mkdir data
COPY data/defaultData.csv ./data
EXPOSE 8080
ENTRYPOINT ["air-quality-mock-endpoint"]