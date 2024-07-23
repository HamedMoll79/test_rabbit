# Build Stage
FROM golang:1.22.4-alpine AS build

RUN apk add --update ca-certificates git make

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY . .

RUN make build-test

# Runtime Stage
FROM alpine:3.20
ENV TZ=UTC

RUN apk add --update --no-cache tzdata ca-certificates bash curl && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    apk del tzdata

COPY --from=build /app/build/test /bin/test
COPY --from=build /app/deploy/test/config.local.yml /bin/config.yml
COPY --from=build /etc/ssl/certs /etc/ssl/certs

# Run
CMD [ "/bin/test"]
