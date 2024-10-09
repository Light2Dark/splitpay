FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/go-starter ./cmd/api 

# run tests in container
FROM build-stage as run-test-stage
RUN go test -v ./...

# deploy app binary into lean image
FROM alpine:edge AS build-release-stage

WORKDIR /

COPY --from=build-stage /bin/go-starter /bin/go-starter
COPY --from=build-stage /app/static /static

# set timezone & install CA certs for HTTPS
RUN apk --no-cache add ca-certificates tzdata

EXPOSE 8080

ENTRYPOINT ["/bin/go-starter"]