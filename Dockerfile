FROM golang:1.20.5-alpine AS dependencies
WORKDIR /dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /build
COPY . .
RUN go build -v -o ./main ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /build/main ./
CMD ["./main"]