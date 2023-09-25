FROM golang:1.20.8-bullseye

WORKDIR /srv/maschera

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /usr/local/bin/maschera

CMD ["/usr/local/bin/maschera"]
