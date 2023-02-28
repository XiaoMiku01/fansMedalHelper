#!/bin/sh
Green="\\033[32m"
Red="\\033[31m"
Plain="\\033[0m"

set -e

case ${MIRRORS} in
"custom")
    # custom
    if [ -z "${CUSTOM_REPO+x}" ]; then
      echo -e "${Red} [ERR] 未配置自定义仓库链接！ ${Plain}"
      exit 1
    else
      echo -e "${Green} [INFO] 使用自定义仓库 ${Plain}"
      git remote set-url origin ${CUSTOM_REPO}
    fi
    ;;
"0")
    # https://github.com/
    echo -e "${Green} [INFO] 使用源-GitHub ${Plain}"
    git remote set-url origin https://github.com/XiaoMiku01/fansMedalHelper.git
    ;;
"1")
    # https://ghproxy.com/
    echo -e "${Green} [INFO] 使用镜像源-GHProxy ${Plain}"
    git remote set-url origin https://ghproxy.com/https://github.com/XiaoMiku01/fansMedalHelper.git
    ;;
"2")
    # http://fastgit.org/
    echo -e "${Green} [INFO] 使用镜像源-FastGIT ${Plain}"
    git remote set-url origin https://hub.fastgit.xyz/XiaoMiku01/fansMedalHelper.git
    ;;
*)
    echo -e "${Green} [INFO] 使用源-GitHub ${Plain}"
    git remote set-url origin https://github.com/XiaoMiku01/fansMedalHelper.git
    ;;
esac

echo -e "${Green} [INFO] 拉取项目更新... ${Plain}"
git config --global --add safe.directory "*"
git pull --no-tags origin master

echo -e "${Green} [INFO] 开始运行... ${Plain}"
python3 main.py