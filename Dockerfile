FROM golang:1.24-alpine as builder  
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o note-tracker cmd/notetracker/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/note-tracker .
CMD ["./note-tracker", "--debug"]
EXPOSE 8080