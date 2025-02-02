FROM golang:1.23.1-alpine
WORKDIR /server
COPY go.mod go.sum ./
COPY . /server
RUN go build -o simple_server .
ENTRYPOINT [ "/server/simple_server" ]