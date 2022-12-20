FROM golang:latest AS dev

RUN mkdir -p /projects/app

WORKDIR /projects/app

COPY go.mod /projects/app/

RUN go mod download

COPY ./*.go /projects/app

RUN pwd

EXPOSE 8080

CMD [ "go", "run", "." ]
