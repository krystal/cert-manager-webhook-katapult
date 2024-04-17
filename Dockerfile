FROM golang:1.22.1-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /work

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -mod=readonly -a -o cert-manager-webhook-katapult .

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /work/cert-manager-webhook-katapult .
ENTRYPOINT ["/cert-manager-webhook-katapult"]
