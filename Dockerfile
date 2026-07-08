FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /scheduler /app/scheduler
COPY web /app/web
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=
EXPOSE 7540
ENTRYPOINT ["/app/scheduler"]