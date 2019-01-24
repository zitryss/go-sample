FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-w'

FROM scratch
COPY --from=builder /app/perfmon /
COPY --from=builder /app/assets/ /assets/
EXPOSE 9000
ENTRYPOINT ["/perfmon"]
