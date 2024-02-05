FROM debian:bookworm-20240130-slim

ARG UID=10001

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
      jq yq && \
    apt-get clean && \
    apt-get autoremove && \
    rm -rf /var/lib/apt/lists/* /tmp/* ~/* && \
    adduser \
      --disabled-password \
      --gecos "" \
      --home "/nonexistent" \
      --shell "/sbin/nologin" \
      --no-create-home \
      --uid "${UID}" \
      kli

USER kli

WORKDIR /app

COPY --chown=$UID kai-kli /app/kai-kli

ENTRYPOINT [ "/app/kai-kli" ]
