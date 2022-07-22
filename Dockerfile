# Container image that runs the code
FROM golang:1.17-alpine

WORKDIR /application

EXPOSE 8085

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./keynetic ./src/*.go

RUN chmod +x ./keynetic

CMD ["./keynetic"]