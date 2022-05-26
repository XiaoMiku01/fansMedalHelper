
from aiohttp import ClientSession, ClientTimeout
import sys
import os
import asyncio

from loguru import logger

sys.path.append(os.path.dirname(
    os.path.dirname(os.path.abspath(__file__))))

logger.remove()
logger.add(sys.stdout, colorize=True,
           format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> <blue> {extra[user]} </blue> <level>{message}</level>", backtrace=True, diagnose=True)


class BiliUser:

    def __init__(self, access_token: str, bannedUIDs: str = ''):
        from .api import BiliApi
        self.mid, self.name = 0, ""
        self.access_key = access_token  # 登录凭证
        self.bannedUIDs = str(bannedUIDs)  # 被禁止的房间ID "1,2,3"
        self.medals = []  # 用户所有勋章
        self.medalsLower20 = []  # 用户所有勋章，等级小于20的

        self.session = ClientSession(timeout=ClientTimeout(total=3))
        self.api = BiliApi(self, self.session)

        self.retryTimes = 0  # 点赞任务重试次数
        self.maxRetryTimes = 10  # 最大重试次数
        self.message = []
        self.errmsg = []

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
        try:
            self.bannedList = list(map(lambda x: int(x if x else 0), self.bannedUIDs.split(',')))
            if self.bannedList:
                self.log.log("WARNING", "已设置黑名单UID: {}".format(' '.join(map(str, self.bannedList))))
        except ValueError:
            self.bannedList = []

    async def getMedals(self):
        '''
        获取用户勋章
        '''
        self.medals.clear()
        self.medalsLower20.clear()
        async for medal in self.api.getFansMedalandRoomID():
            if medal['medal']['target_id'] in self.bannedList:
                continue
            self.medals.append(medal) if medal['room_info']['room_id'] != 0 else None
        [self.medalsLower20.append(medal) for medal in self.medals if medal['medal']['level'] < 20]

    async def likeandShare(self, failedMedals: list = []):
        '''
        点赞 *3 分享 * 5异步执行
        '''
        try:
            self.log.log("INFO", "点赞、分享任务开始....")
            allTasks = []
            if not failedMedals:
                for medal in self.medalsLower20:
                    allTasks.append(self.api.likeInteract(medal['room_info']['room_id']))
                    allTasks.append(self.api.shareRoom(medal['room_info']['room_id']))
            else:
                for medal in failedMedals:
                    allTasks.append(self.api.likeInteract(medal['room_info']['room_id']))
                    allTasks.append(self.api.shareRoom(medal['room_info']['room_id']))
            await asyncio.gather(*allTasks)
            await asyncio.sleep(10)
            await self.getMedals()  # 刷新勋章
            self.log.log("SUCCESS", "点赞、分享任务完成")
            finallyMedals = [medla for medla in self.medalsLower20 if medla['medal']['today_feed'] >= 1200]
            failedMedals = [medla for medla in self.medalsLower20 if medla['medal']['today_feed'] < 1200]
            msg = "20级以下牌子共 {} 个,完成任务 {} 个".format(len(self.medalsLower20), len(finallyMedals))

            self.log.log("INFO", msg)
            self.log.log("WARNING", "失败房间: {}... {}个".format(
                ' '.join([medals['anchor_info']['nick_name'] for medals in failedMedals[:5]]), len(failedMedals)))
            if self.retryTimes > self.maxRetryTimes:
                self.log.log("ERROR", "任务重试次数过多,停止任务")
                return
            if len(finallyMedals) / len(self.medalsLower20) <= 0.9:
                self.log.log("WARNING", "成功率过低,重新执行任务")
                self.retryTimes += 1
                self.log.log("WARNING", "重试次数: {}/{}".format(self.retryTimes, self.maxRetryTimes))
                await self.likeandShare(failedMedals)
            else:
                self.message.append(f"【{self.name}】 " + msg)
                self.errmsg.append(f"【{self.name}】 " + "失败房间: {}... {}个".format(
                    ' '.join([medals['anchor_info']['nick_name'] for medals in failedMedals[:5]]), len(failedMedals)))
        except Exception as e:
            self.log.exception("点赞、分享任务异常")
            self.errmsg.append(f"【{self.name}】 点赞、分享任务异常,请检查日志")

    async def sendDanmaku(self):
        '''
        每日弹幕打卡
        '''
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
                await asyncio.sleep(6)
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
            task = [self.likeandShare(), self.sendDanmaku()]
            await asyncio.wait(task)
        await self.session.close()

    async def sendmsg(self):
        return self.message+self.errmsg
