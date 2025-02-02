FROM golang:1.23.1-alpine 
WORKDIR /weather
COPY go.mod go.sum ./ 
COPY . ./ 
RUN go build -o weather_app . 
ENTRYPOINT ["/weather/weather_app"]