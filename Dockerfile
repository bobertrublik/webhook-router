FROM golang:1.21.5-alpine3.19

RUN apk update && apk upgrade && \
    apk add --no-cache bash build-base

# Define current working directory
WORKDIR /opt/webhook-router

# Download modules to local cache so we can skip re-
# downloading on consecutive docker build commands
COPY go.mod .
COPY go.sum .
RUN go mod download

# Add sources
COPY . .

CMD ["go", "run", "./cmd/webhookd/main.go"]
