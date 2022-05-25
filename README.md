<p align="center">
  <img src="https://s1.ax1x.com/2022/05/24/XPx1tx.png" width="200" height="200" alt="">
</p>
<div align="center">
<h1> 新 B 站粉丝牌助手
</h1>
<p>当前版本：0.1.0</p>
 </div>

**TODO**

-   [x] 每日直播区签到
-   [x] 每日点赞 3 次直播间 （200\*3 亲密度）
-   [x] 每日分享 5 次直播间 （100\*5 亲密度）
-   [x] 每日弹幕打卡 （100 亲密度）
-   [x] 多账号支持
-   [ ] 每日观看 30 分钟
-   [ ] 微信推送通知

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
    - access_key: XXXXXX # 注意冒号后的空格 否则会读取失败
      banned_uid: 789,100 # 黑名单UID 同上,填了后将不会打卡，点赞，分享
    - access_key:
      banned_uid:
    # 多用户以上格式添加
    # 井号后为注释 井号前后必须有空格
CRON: # 0 0 * * *
# 这里是 cron 表达式, 第一个参数是分钟, 第二个参数是小时
# 例如每天凌晨0点0分执行一次为 0 0 * * *
# 如果不填,则不使用内置定时器,填写正确后要保持该进程一直运行
```

请务必严格填写，否则程序将读取失败，可以在这里 [YAML、YML 在线编辑器(格式化校验)-BeJSON.com](https://www.bejson.com/validators/yaml_editor/) 验证你填的 yaml 是否正确

**注意！！ 如果填写 28 个或者很多需要分享的 UID，最好是设置每天 0 点运行，不然会出现还没运行完就到第二天刷新任务了**

> 运行主程序

```shell
python main.py
```

> 效果图

[![XiifQP.md.png](https://s1.ax1x.com/2022/05/24/XiifQP.md.png)](https://imgtu.com/i/XiifQP)

---

### 注意事项

---

### 更新日志

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
