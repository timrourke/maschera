FROM golang:1.20.8-bullseye

WORKDIR /srv/maschera

RUN go install github.com/cosmtrek/air@v1.45.0

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

CMD ["air"]
