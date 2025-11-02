FROM golang:1.25.3-bookworm as go

FROM node:25-bookworm
COPY --from=go /usr/local/go /usr/local/go
ENV PATH "/usr/local/go/bin:$PATH"
ENV GOPATH /go
ENV GOCACHE /go/cache
RUN apt install tzdata -y

WORKDIR /app
CMD [ "sleep", "infinity" ]
