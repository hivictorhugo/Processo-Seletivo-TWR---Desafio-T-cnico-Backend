FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/server ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/server /bin/server
EXPOSE 8080
CMD ["/bin/server"]