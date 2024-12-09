FROM golang:1.23-alpine AS builder

RUN apk --no-cache add ca-certificates

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/src/app/dist/ ./...

FROM scratch AS runtime

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/dist/hcloud-upload-image /bin/hcloud-upload-image

ENTRYPOINT ["/bin/hcloud-upload-image"]
