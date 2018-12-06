# build app
FROM golang AS build-env

RUN go get -u github.com/swaggo/swag/cmd/swag

ADD . /app

WORKDIR /app

RUN swag init
RUN go build -o app

# safe image
FROM debian

ENV TZ=Asia/Shanghai

COPY --from=build-env /app/app /usr/bin/app

EXPOSE 1324

CMD ["app"]
