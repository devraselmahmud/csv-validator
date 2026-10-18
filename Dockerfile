FROM golang:1.27.2 AS build
WORKDIR /src
COPY go.mod ./
COPY internal/ ./internal/
COPY cmd/ ./cmd/
RUN go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /csv-validator ./cmd/csv-validator

FROM scratch
COPY --from=build /csv-validator /csv-validator
USER 65532:65532
WORKDIR /data
ENTRYPOINT ["/csv-validator"]
