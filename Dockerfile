# test & build app
FROM golang AS build-env

ADD . /app

WORKDIR /app

RUN go test ./...
RUN go build -o api cmd/api

# safe image
FROM muninn/debian

COPY --from=build-env /app/api /usr/bin/api

EXPOSE 80

HEALTHCHECK CMD curl -f http://localhost/status || exit 1

CMD ["api"]
