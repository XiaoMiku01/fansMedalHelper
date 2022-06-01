FROM python:3.9-slim
ENV TZ="Asia/Shanghai"
RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app/fansMedalHelper

COPY . .

RUN  pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple


CMD [ "python", "index.py" ]