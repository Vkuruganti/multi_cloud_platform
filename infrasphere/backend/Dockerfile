FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/api ./cmd/api && go build -o /out/worker ./cmd/worker

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker
EXPOSE 8080
CMD ["/app/api"]

