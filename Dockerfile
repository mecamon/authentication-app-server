
FROM golang:1.17-alpine AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/test ./cmd/web/*.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /usr/src/app/bin /go/bin
EXPOSE 8080
ENTRYPOINT /go/bin/test