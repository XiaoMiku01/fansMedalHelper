FROM python:3.9-slim
ENV TZ="Asia/Shanghai"

WORKDIR /app/fansMedalHelper

COPY . .

RUN  pip3 install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple


CMD [ "python", "main.py" ]
