FROM golang:1.17-alpine as builder

RUN apk add --no-cache make

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN make build

FROM scratch

COPY --from=builder /src/build/tp-link-exporter /usr/local/bin/

ENTRYPOINT ["tp-link-exporter"]
