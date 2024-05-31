FROM golang:alpine AS builder
WORKDIR /application
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-rest-api

FROM scratch
ENV MYSQL_USER=
ENV MYSQL_PASS=
ENV MYSQL_SCHEMA=
ENV MYSQL_ENDPOINT=
ENV GOAPI_ENDPOINT=
COPY --from=builder /go-rest-api /go-rest-api
ENTRYPOINT ["/go-rest-api"]