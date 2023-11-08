FROM golang:1.19

WORKDIR /app

COPY server /app/

COPY client /app/

COPY script.sh /app/
WORKDIR /app/client

RUN go mod download

WORKDIR /app/server

RUN go mod download

WORKDIR /app

EXPOSE 9000
EXPOSE 8080

# Run
CMD ["/script.sh"]