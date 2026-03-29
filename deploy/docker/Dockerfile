FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/sub-hub ./cmd/sub-hub

FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=build /out/sub-hub /app/sub-hub
COPY api/openapi.yaml /app/api/openapi.yaml

ENV OPENAPI_PATH=/app/api/openapi.yaml
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/sub-hub"]

