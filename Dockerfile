FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stuber ./cmd/stuber

FROM scratch

COPY --from=builder /app/stuber /stuber

ENTRYPOINT ["/stuber"]
CMD ["up", "-f", "config.yaml"]
