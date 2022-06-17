# Start from golang base image
FROM golang:1.18-alpine as builder

WORKDIR /app/cacheImage

COPY . .

ENV GO111MODULE=on

RUN go build -o cacheImage -ldflags='-s -w' ./cmd/cli/main.go

# Add Maintainer info
FROM golang:1.18-alpine

# creates working directory for program
WORKDIR /app/cacheImage

COPY --from=builder /app/cacheImage/cacheImage .

EXPOSE ${PORT}

CMD [ "./cacheImage" ]
