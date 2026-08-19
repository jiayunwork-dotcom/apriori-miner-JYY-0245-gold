FROM golang:1.21-alpine

ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build ./...

CMD ["bash"]
