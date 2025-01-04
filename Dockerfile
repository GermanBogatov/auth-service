FROM golang:1.23-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go clean --modcache && \
    go mod download && \
    go mod verify

COPY . .
RUN go build -v -o app cmd/application/main.go

FROM golang:1.23-alpine

WORKDIR /application

COPY --from=builder /build/app /application

EXPOSE 8080
CMD ./app