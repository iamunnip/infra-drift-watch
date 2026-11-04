FROM golang:1.27.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -trimpath -o /out/infra-drift-watch ./cmd/infra-drift-watch

FROM alpine:3.24.3 AS tools
RUN apk add --no-cache ca-certificates curl unzip diffutils
COPY scripts/install-tools.sh /install-tools.sh
RUN SKIP_TOFU=1 /install-tools.sh /tools

FROM alpine:3.24.3
RUN apk add --no-cache ca-certificates diffutils \
    && mkdir -p /work /state \
    && chown 65532:65532 /work /state
COPY --from=build /out/infra-drift-watch /usr/local/bin/infra-drift-watch
COPY --from=tools /tools/ /usr/local/bin/
COPY --from=ghcr.io/opentofu/opentofu:1.13.5-minimal /usr/local/bin/tofu /usr/local/bin/tofu
COPY LICENSE /LICENSE
COPY THIRD_PARTY.md /THIRD_PARTY.md
USER 65532:65532
WORKDIR /work
ENTRYPOINT ["/usr/local/bin/infra-drift-watch"]
CMD ["--help"]
