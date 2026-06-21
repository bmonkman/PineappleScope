# Build stage: compile the CGO sqlite binary natively for linux.
# (Building inside Docker on a Linux runner avoids the cross-compilation pain
# that previously required xgo.)
FROM golang:1.21-bookworm AS builder

WORKDIR /src

# Cache module downloads across builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# VERSION (e.g. the git SHA) is baked in so asset URLs bust caches per build.
ARG VERSION=dev
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o /out/pineapplescope ./cmd/pineapplescope

# Runtime stage: slim image with just the binary, templates/assets, and certs.
FROM debian:bookworm-slim

ENV TZ=America/Vancouver
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/pineapplescope /app/pineapplescope
COPY resources/ /app/resources/

EXPOSE 1111

# DB lives on a volume so it survives container recreation/updates.
ENV DBFILE=/var/db/pineapplescope.db
VOLUME /var/db/

ENTRYPOINT ["/app/pineapplescope"]
