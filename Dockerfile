FROM golang

COPY ./src /go/src/github.com/jeckel/flashcrowd/src
WORKDIR /go/src/github.com/jeckel/flashcrowd/src

RUN go get ./
RUN go build

CMD flashcrowd
