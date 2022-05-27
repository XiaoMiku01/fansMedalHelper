FROM ubuntu:latest
#维护镜像的用户信息.
MAINTAINER dovela
#镜像操作指令安装apache软件
RUN apt update \
	&& apt install --no-install-recommends python3-pip git -y \
	&& git clone https://ghproxy.com/https://github.com/dovela/fansMedalHelper \
	&& pip3 install -r /fansMedalHelper/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple \
	&& apt purge --autoremove git -y \
	&& apt purge python3-pip -y \
	&& rm -rf /var/lib/apt/lists/*
WORKDIR /fansMedalHelper

CMD [ "python3","main.py"]