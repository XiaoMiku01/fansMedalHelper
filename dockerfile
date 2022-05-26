FROM debian:latest
#维护镜像的用户信息.
MAINTAINER dovela
#镜像操作指令安装apache软件
RUN apt update \
	&& apt install python3-pip git -y \
	&& git clone https://github.com/XiaoMiku01/fansMedalHelper.git \
	&& pip3 install -r /fansMedalHelper/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
WORKDIR /fansMedalHelper

CMD [ "python3","main.py"]