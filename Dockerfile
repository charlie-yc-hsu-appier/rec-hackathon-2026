FROM golang:1.25-alpine AS build-env

ENV GOPRIVATE=github.com/plaxieappier

WORKDIR /rec-vendor-api

RUN apk add build-base git ca-certificates openssh openssl
RUN mkdir -pm 0600 /root/.ssh \
    && touch /root/.ssh/known_hosts \
    && ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6 \
    && go install go.uber.org/mock/mockgen@v0.4.0

COPY .gitconfig /root/.gitconfig
COPY go.* ./

RUN --mount=type=ssh,id=ai-rec-common go mod download

COPY ./cmd ./cmd/
COPY ./internal ./internal/

RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server.app ./cmd/rec-vendor-api/server.go


FROM scratch
WORKDIR /srv
COPY --from=build-env /rec-vendor-api/bin/server.app /srv
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/srv/server.app"]
