FROM golang:1.18.1

WORKDIR /app

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /tcpserver

EXPOSE 1234

CMD ["/tcpserver", "-host", "0.0.0.0"]

