FROM golang AS build

RUN apt update && apt upgrade -y

COPY go.mod go.sum /build/
COPY src /build/src

RUN cd /build/src && \ 
    go install && \
    go mod download && \ 
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o inciprone -ldflags="-extldflags=-static" 

FROM scratch
COPY --from=build /build/src/inciprone .
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/

ENTRYPOINT ["./inciprone"]