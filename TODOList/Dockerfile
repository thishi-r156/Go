FROM golang:1.23.1-alpine 
WORKDIR /todolist
COPY go.mod go.sum ./
COPY . ./
RUN go build -o todoapp . 
ENTRYPOINT ["/todolist/todoapp"] 