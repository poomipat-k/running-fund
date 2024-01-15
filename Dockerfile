# Build Stage

FROM golang:1.21.5-alpine3.19 as buildStage

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

RUN env GOOS=linux CGO_ENABLED=0 go build -o /myApp ./cmd/

# Deploy Stage

FROM alpine:latest

WORKDIR /

COPY --from=buildStage /myApp /app

COPY --from=buildStage /app/resources /home

EXPOSE 8080

ENTRYPOINT [ "/app" ]