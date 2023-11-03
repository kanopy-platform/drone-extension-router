FROM golang:1.21 as build
ARG VERSION="0.0.0"
ARG GIT_COMMIT
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-X 'github.com/kanopy-platform/drone-extension-router/internal/version.version=${VERSION}' -X 'github.com/kanopy-platform/drone-extension-router/internal/version.gitCommit=${GIT_COMMIT}'" -o /go/bin/app

FROM debian:bookworm-slim
RUN apt-get update && apt-get install --yes ca-certificates
RUN groupadd -r app && useradd --no-log-init -r -g app app
USER app
COPY --from=build /go/bin/app /
ENV APP_ADDR ":8080"
EXPOSE 8080
ENTRYPOINT ["/app"]
