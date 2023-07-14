#FROM golang:1.19 as builder
#COPY go.mod go.sum /go/src/ticket-expert/
#WORKDIR /go/src/ticket-expert
#RUN go mod download
#COPY . /go/src/ticket-expert
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/ticket-expert ticket-expert
#
#FROM alpine
#RUN apk add --no-cache ca-certificates && update-ca-certificates
#COPY --from=builder /go/src/ticket-expert/build/ticket-expert /usr/bin/ticket-expert
#EXPOSE 8080 8080
#ENTRYPOINT ["/usr/bin/ticket-expert"]

# syntax=docker/dockerfile:1

FROM golang:1.19
#ENV GOPATH=/ticket-expert/main
# Set destination for COPY
WORKDIR /ticket-expert

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./
RUN pwd
RUN ls
RUN echo "~~~~~~~~"
RUN printenv
WORKDIR /
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ticket-expert/main

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]