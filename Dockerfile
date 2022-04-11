FROM golang:1.18-alpine
RUN apk add  --no-cache ffmpeg
COPY ./app /app
WORKDIR /app
RUN go mod download
RUN go build -o app
CMD ["go", "run", "main.go"]