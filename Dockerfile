FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o note-tracker ./cmd/notetracker/main.go


FROM alpine:latest AS runner
WORKDIR /root/
COPY --from=builder /app/note-tracker .
EXPOSE 8080
CMD ["./note-tracker", "--port=8080", "--host=0.0.0.0", "--debug"]