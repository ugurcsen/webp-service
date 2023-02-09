FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/main .

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY --chown=root:root --chmod=755 ./bin/ /app/bin/
ENV PORT 8080
ENV FILE_SIZE_LIMIT 5000
CMD ["/app/main"]