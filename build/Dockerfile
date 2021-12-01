FROM python:3.8.12-slim-buster AS runner
RUN apt-get update && apt-get install -y wget xz-utils libgl1-mesa-glx
RUN wget -nv -O- https://download.calibre-ebook.com/linux-installer.sh | sh /dev/stdin

FROM runner

ARG BUILD_VERSION=0.1

LABEL org.opencontainers.image.version=$BUILD_VERSION \
      org.opencontainers.image.authors="airenass@gmail.com" \
      name="get-text" \
      maintainer="airenass@gmail.com" \
      vendor="airenass@gmail.com" \
      version=$BUILD_VERSION \
      release=$BUILD_VERSION \
      summary="This image is a wrapper for calibre's ebook-converter" \
      description="This image is used to host the wrapper service for calibre's ebook-converter" 

WORKDIR /app

RUN adduser --system --group --no-create-home app

RUN mkdir -p /app && chown -R app:app /app   
RUN mkdir -p /tmp && chown -R app:app /tmp  

COPY data/get-text /app/
COPY config.yaml /app/

RUN chown app:app /app/* /app

USER app

ENTRYPOINT ["./get-text"]
