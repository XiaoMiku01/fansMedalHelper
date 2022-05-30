
from aiohttp import ClientSession, ClientTimeout
import sys
import os
import asyncio
import uuid
from loguru import logger

sys.path.append(os.path.dirname(
    os.path.dirname(os.path.abspath(__file__))))

logger.remove()
logger.add(sys.stdout, colorize=True,
           format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> <blue> {extra[user]} </blue> <level>{message}</level>", backtrace=True, diagnose=True)


class BiliUser:

    def __init__(self, access_token: str, whiteUIDs: str = '', bannedUIDs: str = '', config: dict = {}):
        from .api import BiliApi
        self.mid, self.name = 0, ""
        self.access_key = access_token  # 登录凭证
        try:
            self.whiteList = list(map(lambda x: int(x if x else 0), str(whiteUIDs).split(',')))  # 白名单UID
            self.bannedList = list(map(lambda x: int(x if x else 0), str(bannedUIDs).split(',')))  # 黑名单
        except ValueError:
            raise ValueError("白名单或黑名单格式错误")
        self.config = config
        self.medals = []  # 用户所有勋章
        self.medalsLower20 = []  # 用户所有勋章，等级小于20的

        self.session = ClientSession(timeout=ClientTimeout(total=3))
        self.api = BiliApi(self, self.session)

        self.retryTimes = 0  # 点赞任务重试次数
        self.maxRetryTimes = 10  # 最大重试次数
        self.message = []
        self.errmsg = ["错误日志："]
        self.uuids = [str(uuid.uuid4()) for _ in range(2)]

    async def loginVerify(self) -> bool:
        '''
        登录验证
        '''
        loginInfo = await self.api.loginVerift()
        self.mid, self.name = loginInfo['mid'], loginInfo['name']
        self.log = logger.bind(user=self.name)
        if loginInfo['mid'] == 0:
            self.isLogin = False
            return False
        self.log.log("SUCCESS", str(loginInfo['mid']) + " 登录成功")

        self.isLogin = True
        return True

    async def doSign(self):
        try:
            signInfo = await self.api.doSign()
            self.log.log("SUCCESS", "签到成功,本月签到次数: {}/{}".format(signInfo['hadSignDays'], signInfo['allDays']))
            self.message.append(f"【{self.name}】 签到成功,本月签到次数: {signInfo['hadSignDays']}/{signInfo['allDays']}")
        except Exception as e:
            self.log.log("ERROR", e)
            self.errmsg.append(f"【{self.name}】" + str(e))
        userInfo = await self.api.getUserInfo()
        self.log.log("INFO", "当前用户UL等级: {} ,还差 {} 经验升级".format(userInfo['exp']['user_level'], userInfo['exp']['unext']))
        self.message.append(
            f"【{self.name}】 UL等级: {userInfo['exp']['user_level']} ,还差 {userInfo['exp']['unext']} 经验升级")

    async def getMedals(self):
        '''
        获取用户勋章
        '''
        self.medals.clear()
        self.medalsLower20.clear()
        async for medal in self.api.getFansMedalandRoomID():
            if self.whiteList == [0]:
                if medal['medal']['target_id'] in self.bannedList:
                    self.log.warning(f"{medal['anchor_info']['nick_name']} 在黑名单中，已过滤")
                    continue
                self.medals.append(medal) if medal['room_info']['room_id'] != 0 else ...
            else:
                if medal['medal']['target_id'] in self.whiteList:
                    self.medals.append(medal) if medal['room_info']['room_id'] != 0 else ...
                    self.log.success(f"{medal['anchor_info']['nick_name']} 在白名单中，加入任务")
        [self.medalsLower20.append(medal) for medal in self.medals if medal['medal']['level'] < 20]

    async def asynclikeandShare(self, failedMedals: list = []):
        '''
        点赞 *3 分享 * 5 
        '''
        if self.config['LIKE_CD'] == 0:
            self.log.log("INFO", "点赞任务已关闭")
        elif self.config['SHARE_CD'] == 0:
            self.log.log("INFO", "分享任务已关闭")
        if self.config['LIKE_CD'] == 0 and self.config['SHARE_CD'] == 0:
            return
        if not self.config['ASYNC']:
            self.log.log("INFO", "同步点赞、分享任务开始....")
            for index, medal in enumerate(self.medalsLower20):
                tasks = []
                tasks.append(self.api.likeInteract(medal['room_info']['room_id'])) if self.config['LIKE_CD'] else ...
                tasks.append(self.api.shareRoom(medal['room_info']['room_id'])) if self.config['SHARE_CD'] else ...
                await asyncio.gather(*tasks)
                self.log.log("SUCCESS", f"{medal['anchor_info']['nick_name']} 点赞,分享成功 {index+1}/{len(self.medalsLower20)}")
                await asyncio.sleep(max(self.config['LIKE_CD'], self.config['SHARE_CD']))
            return
        try:
            self.log.log("INFO", "异步点赞、分享任务开始....")
            allTasks = []
            if not failedMedals:
                failedMedals = self.medalsLower20
            for medal in failedMedals:
                allTasks.append(self.api.likeInteract(medal['room_info']['room_id'])) if self.config['LIKE_CD'] else ...
                allTasks.append(self.api.shareRoom(medal['room_info']['room_id'])) if self.config['SHARE_CD'] else ...
            await asyncio.gather(*allTasks)
            await asyncio.sleep(10)
            await self.getMedals()  # 刷新勋章
            self.log.log("SUCCESS", "点赞、分享任务完成")
            finallyMedals = [medla for medla in self.medalsLower20 if medal['medal']['today_feed'] >= 1200]
            midMedals = [medla for medla in self.medalsLower20 if medal['medal']['today_feed'] >= 1100]
            failedMedals = [medla for medla in self.medalsLower20 if medal['medal']['today_feed'] < 1100]
            msg = "20级以下牌子共 {} 个,完成任务 {} 个亲密度大于1100, {} 个亲密度大于1200".format(
                len(self.medalsLower20), len(midMedals), len(finallyMedals))

            self.log.log("INFO", msg)
            self.log.log("WARNING", "小于1100或失败房间: {}... {}个".format(
                ' '.join([medals['anchor_info']['nick_name'] for medals in failedMedals[:5]]), len(failedMedals)))
            if self.retryTimes > self.maxRetryTimes:
                self.log.log("ERROR", "任务重试次数过多,停止任务")
                return
            if len(finallyMedals) / len(self.medalsLower20) <= 0.9:
                self.log.log("WARNING", "成功率过低,重新执行任务")
                self.retryTimes += 1
                self.log.log("WARNING", "重试次数: {}/{}".format(self.retryTimes, self.maxRetryTimes))
                await self.asynclikeandShare(failedMedals)
            else:
                self.message.append(f"【{self.name}】 " + msg)
                self.errmsg.append(f"【{self.name}】 " + "小于1100或失败房间: {}... {}个".format(
                    ' '.join([medals['anchor_info']['nick_name'] for medals in failedMedals[:5]]), len(failedMedals)))
        except Exception as e:
            self.log.exception("点赞、分享任务异常")
            self.errmsg.append(f"【{self.name}】 点赞、分享任务异常,请检查日志")

    async def sendDanmaku(self):
        '''
        每日弹幕打卡
        '''
        if not self.config['DANMAKU_CD']:
            self.log.log("INFO", "弹幕任务关闭")
            return
        self.log.log("INFO", "弹幕打卡任务开始....(预计 {} 秒完成)".format(len(self.medals) * 6))
        n = 0
        for medal in self.medals:
            try:
                danmaku = await self.api.sendDanmaku(medal['room_info']['room_id'])
                n += 1
                self.log.log(
                    "DEBUG", "{} 房间弹幕打卡成功: {} ({}/{})".format(medal['anchor_info']['nick_name'], danmaku, n, len(self.medals)))
            except Exception as e:
                self.log.log("ERROR", "{} 房间弹幕打卡失败: {}".format(medal['anchor_info']['nick_name'], e))
                self.errmsg.append(f"【{self.name}】 {medal['anchor_info']['nick_name']} 房间弹幕打卡失败: {str(e)}")
            finally:
                await asyncio.sleep(self.config['DANMAKU_CD'])
        self.log.log("SUCCESS", "弹幕打卡任务完成")
        self.message.append(f"【{self.name}】 弹幕打卡任务完成 {n}/{len(self.medals)}")

    async def init(self):
        if not await self.loginVerify():
            self.log.log("ERROR", "登录失败")
            self.errmsg.append("登录失败")
            await self.session.close()
        else:
            await self.doSign()
            await self.getMedals()

    async def start(self):
        if self.isLogin:
            task = [self.asynclikeandShare(), self.sendDanmaku(), self.watchinglive()]
            await asyncio.wait(task)
        # await self.session.close()

    async def sendmsg(self):
        if not self.isLogin:
            await self.getMedals()
            await self.session.close()
            return self.message+self.errmsg
        nameList1, nameList2, nameList3, nameList4 = [], [], [], []
        for medal in self.medalsLower20:
            today_feed = medal['medal']['today_feed']
            nick_name = medal['anchor_info']['nick_name']
            if today_feed >= 1300:
                nameList1.append(nick_name)
            elif 1200 <= today_feed < 1300:
                nameList2.append(nick_name)
            elif 1100 <= today_feed < 1200:
                nameList3.append(nick_name)
            elif today_feed < 1100:
                nameList4.append(nick_name)
        self.message.append(f"【{self.name}】 今日亲密度获取情况如下（20级以下）：")

        for l, n in zip([nameList1, nameList2, nameList3, nameList4], ["【1300及以上】", "【1200至1300】", "【1100至1200】", "【1100以下】"]):
            if len(l) > 0:
                self.message.append(f"{n}" + ' '.join(l[:5]) + f"{'等' if len(l) > 5 else ''}" + f' {len(l)}个')
        await self.session.close()
        return self.message+self.errmsg

    async def watchinglive(self):
        if not self.config['WATCHINGLIVE']:
            self.log.log("INFO", "每日30分钟任务关闭")
            return
        self.log.log("INFO", "每日30分钟任务开始")
        heartNum = 0
        while True:
            tasks = []
            for medal in self.medalsLower20:
                tasks.append(self.api.heartbeat(medal['room_info']['room_id'], medal['medal']['target_id']))
            await asyncio.wait(tasks)
            heartNum += 1
            self.log.log(
                "INFO", f"{' '.join([medal['anchor_info']['nick_name'] for medal in self.medalsLower20[:5]])} 等共 {len(self.medalsLower20)} 个房间的第{heartNum}次心跳包已发送（{heartNum}/{30}）")
            await asyncio.sleep(60)
            if heartNum >= 30:
                break
        self.log.log("SUCCESS", "每日30分钟任务完成")
