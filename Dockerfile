FROM python:3.9-alpine
ENV TZ="Asia/Shanghai"

WORKDIR /tmp

RUN apk add --no-cache git \
    && git config --global --add safe.directory "*" \
    && git clone https://github.com/XiaoMiku01/fansMedalHelper /app/fansMedalHelper \
    && pip install --no-cache-dir -r /app/fansMedalHelper/requirements.txt \
    && rm -rf /tmp/*

WORKDIR /app/fansMedalHelper

ENTRYPOINT ["/bin/sh","/app/fansMedalHelper/entrypoint.sh"]