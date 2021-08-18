FROM golang:1.14-alpine AS build-env
RUN apk add librdkafka-dev build-base git openssh-client

ARG SSH_PRIVATE_KEY

WORKDIR /usr/src/app

# Allow pulling private modules by adding private ssh key and adding github.com to known hosts
RUN git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"

RUN mkdir -p /root/.ssh
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

RUN md5sum /root/.ssh/id_rsa

RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# We also need to tell go which repositories are considered private
RUN go env -w GOPRIVATE=github.com/ridebeam

# trying to cache modules independent of source code changes
COPY go.mod go.sum  ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p bin
RUN go build -tags musl -o bin ./...
RUN go test -tags musl ./...

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
ENV KAFKA_SSL_CA /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=build-env /usr/src/app/bin .

EXPOSE 80
CMD ["./change-me-service-name"]
