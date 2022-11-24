FROM golang:1.18.7-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY css ./css
COPY html ./html
COPY js ./js
COPY src ./src
COPY blog ./blog
COPY *.go ./

RUN go build main.go

ENTRYPOINT [ "./main" ]
