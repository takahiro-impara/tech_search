FROM golang:1.18.1-alpine

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

EXPOSE 8080

RUN apk update && apk add git
COPY go.mod go.sum main.go ./
RUN go mod download

CMD ["go", "run", "main.go"]