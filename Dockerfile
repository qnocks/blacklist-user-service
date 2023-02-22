FROM golang:1.19-bullseye

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o app ./cmd/app/main.go

CMD ["./app"]
