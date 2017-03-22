FROM alpine:latest

COPY ./etc/timezone /etc/timezone

COPY ./etc/localtime /etc/localtime

COPY ./main /bin/kk-top

RUN chmod +x /bin/kk-top

COPY ./config /config

COPY ./app.ini /app.ini

ENV KK_ENV_CONFIG /config/env.ini

VOLUME /config

CMD kk-top $KK_ENV_CONFIG

