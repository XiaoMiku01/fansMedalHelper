<p align="center">
  <img src="https://s1.ax1x.com/2022/05/24/XPx1tx.png" width="200" height="200" alt="">
</p>
<div align="center">
<h1> 新 B 站粉丝牌助手
</h1>

<p>当前版本：0.3.0</p>

 </div>

**TODO**

-   [x] 每日直播区签到
-   [x] 每日点赞 3 次直播间 （200\*3 亲密度）
-   [x] 每日分享 5 次直播间 （100\*5 亲密度）
-   [x] 每日弹幕打卡 （100 亲密度）
-   [x] 每日观看 30 分钟 （100 亲密度）
-   [x] 多账号支持
-   [x] 微信推送通知

<small>ps: 新版 B 站粉丝牌的亲密度每一个牌子都将单独计算  </small>

---

### 使用说明

##### 环境需求：Python 版本大于 3.8

> 克隆本项目 安装依赖

```shell
git clone https://github.com/XiaoMiku01/fansMedalHelper.git
cd fansMedalHelper
pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
```

> 获取 B 站账号的 access_key

下载获取工具 [Release B 站 access_key 获取工具 · XiaoMiku01/fansMedalHelper (github.com)](https://github.com/XiaoMiku01/fansMedalHelper/releases/tag/logintool)

双击打开，扫码登录，会得到 `access_key` 即可

> 填写配置文件 users.yaml

```shell
vim users.yaml
```

```yaml
USERS:
    - access_key: XXXXX # 注意冒号后的空格 否则会读取失败 英文冒号
      white_uid: 0 # 白名单用户ID, 可以是多个用户ID, 以逗号分隔,填写后只会打卡这些用户,黑名单失效，不用就填0
      banned_uid: 0 # 黑名单UID 同上,填了后将不会打卡，点赞，分享 用英文逗号分隔 不填则不限制,两个都填0则不限制,打卡所有直播间
    - access_key:
      white_uid: 0
      banned_uid: 0
    # 注意对齐
    # 多用户以上格式添加
    # 井号后为注释 井号前后必须有空格！井号前后必须有空格！井号前后必须有空格！
    # 冒号后面也要有空格！冒号前面也要有空格！冒号前面也要有空格！
    # 英文冒号,英文逗号！英文逗号！英文逗号！
CRON: # 0 0 * * *
# 这里是 cron 表达式, 第一个参数是分钟, 第二个参数是小时
# 例如每天凌晨0点0分执行一次为 0 0 * * *
# 如果不填,则不使用内置定时器,填写正确后要保持该进程一直运行

SENDKEY: # Server酱微信推送 可选 获取地址：https://sct.ftqq.com/

#########以下为自定义配置#########

ASYNC: 1 # 异步执行,默认异步执行,设置为0则同步执行,开启异步后,将不支持设置点赞和分享CD时间

LIKE_CD: 2 # 点赞间隔时间,单位秒,默认2秒,仅为同步时生效,设置为0则不点赞

SHARE_CD: 5 # 分享间隔时间,单位秒,默认5秒,仅为同步时生效,设置为0则不分享

DANMAKU_CD: 6 # 弹幕间隔时间,单位秒,默认6秒,设置为0则不发弹幕打卡,只能同步打卡

WATCHINGLIVE: 1 # 是否完成每日三十分钟看直播任务，默认开启，设置为0则关闭

# 说明：
# 本项目中的异步执行指的是：同时点赞或者分享所有直播间，速度非常快，但缺点就是可能会被B站吞掉亲密度，所以建议粉丝牌较少的用户开启异步执行
# 粉丝牌数大于30的用户建议使用同步，会更加稳定。缺点就是速度比较慢，但是可以设置点赞和分享的CD时间，避免被B站吞掉亲密度
# 多用户之间依然是异步，不受配置影响
```

请务必严格填写，否则程序将读取失败，可以在这里 [YAML、YML 在线编辑器(格式化校验)-BeJSON.com](https://www.bejson.com/validators/yaml_editor/) 验证你填的 yaml 是否正确

> 运行主程序

```shell
python main.py
```

> 效果图

[![XiifQP.md.png](https://s1.ax1x.com/2022/05/24/XiifQP.md.png)](https://imgtu.com/i/XiifQP)

---

### 已知问题

-   异步执行太快导致部分点赞分享被 B 站吞了。事实上大部分都是 1100 以上

---

### 更新日志

- 2022-5-26

  -   自动完成每日观看30分钟任务

  -   高度自定义的用户配置（详细看`users.yaml`文件）
  -   更加详细的微信推送通知（今日亲密度具体获取情况）

-   2022-5-26

    -   增加异常重试处理
    -   更加详细的报错日志
    -   增加了微信通知功能

-   2022-5-25

    -   B 站取消了分享的 10 分钟 CD，目前已改为异步执行

    -   增加了黑名单设置

    -   增加了自动分享 28 个直播的设置
    -   修复了，粉丝牌过多导致获取不全情况
    -   修复了，粉丝牌过多导致点赞不完全的情况
    -   自动切换运行目录

---

### 赞助

![](http://i0.hdslb.com/bfs/album/c267037c9513b8e44bc6ec95dbf772ff0439dce6.jpg)
