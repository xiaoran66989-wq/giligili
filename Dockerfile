FROM golang:alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/giligili .

FROM alpine:latest

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app
WORKDIR /app
COPY --from=build --chown=app:app /out/giligili /usr/local/bin/giligili
COPY --chown=app:app conf /app/conf
COPY --chown=app:app swagger /app/swagger
USER app
EXPOSE 3000
ENTRYPOINT ["/usr/local/bin/giligili"]
