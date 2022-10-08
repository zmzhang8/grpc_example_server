FROM debian:11-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workdir

RUN useradd service-user \
    && chown -R service-user /workdir

USER service-user

COPY build/server .

ENTRYPOINT ["/workdir/server"]