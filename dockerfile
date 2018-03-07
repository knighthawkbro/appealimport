FROM golang:latest

LABEL maintainer="knighthawk0192@gmail.com"

WORKDIR /go/src/appealimport
COPY . .

# ENV GOPATH /go
# ENV PATH $PATH:/go/bin

RUN go get -d -v ./...
RUN go build -o ./main .

# EXPOSE 8000

# ENTRYPOINT [ "go run", "./appeal.go" ]
CMD ["./main"]