FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/main .

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY ./bin/ /app/bin/
CMD ["/app/main"]