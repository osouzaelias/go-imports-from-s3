FROM golang:1.20
LABEL maintainer="osouzaelias@gmail.com"

COPY . /app

WORKDIR /app/cmd

RUN go build -o go-import-from-s3

CMD ["./go-import-from-s3"]