# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY client ./
COPY server ./
COPY script.sh ./

WORKDIR /app/client

RUN go mod download

WORKDIR /app/server

RUN go mod download

WORKDIR /app
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

# Build

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080
EXPOSE 9000


# Run
CMD ["./script.sh"]