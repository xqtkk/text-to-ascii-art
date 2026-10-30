FROM golang:1.27-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ascii-art-web ./cmd

FROM alpine:3.22
LABEL org.opencontainers.image.title="ASCII Art Web" \
      org.opencontainers.image.description="A web app that converts text into ASCII banner art" \
      org.opencontainers.image.source="https://github.com/xqtkk/text-to-ascii-art"

RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/ascii-art-web /usr/local/bin/ascii-art-web
USER app
EXPOSE 4000
ENTRYPOINT ["ascii-art-web"]
