
from aiohttp import ClientSession
import sys
import os
import asyncio

from loguru import logger

sys.path.append(os.path.dirname(
    os.path.dirname(os.path.abspath(__file__))))

logger.remove()
logger.add(sys.stdout, colorize=True,
           format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> <blue> {extra[user]} </blue> <level>{message}</level>")


class BiliUser:

    def __init__(self, access_token: str, needShareUIDs: str = ''):
        from .api import BiliApi

        self.access_key = access_token  # 登录凭证
        self.needShareUIDs = str(needShareUIDs)  # 需要分享的房间ID "1,2,3"
        self.medals = []  # 用户所有勋章
        self.medalsLower20 = []  # 用户所有勋章，等级小于20的
        self.medalsNeedShare = []  # 用户所有勋章，需要分享的 最多28个

        self.session = ClientSession()
        self.api = BiliApi(self, self.session)

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
        except Exception as e:
            self.log.log("ERROR", e)
        userInfo = await self.api.getUserInfo()
        self.log.log("INFO", "当前用户UL等级: {} ,还差 {} 经验升级".format(userInfo['exp']['user_level'], userInfo['exp']['unext']))

    async def getMedals(self):
        '''
        获取用户勋章
        '''
        self.medals.clear()
        self.medalsLower20.clear()
        async for medal in self.api.getFansMedalandRoomID():
            self.medals.append(medal) if medal['room_info']['room_id'] != 0 else None
        [self.medalsLower20.append(medal) for medal in self.medals if medal['medal']['level'] < 20]
        if self.needShareUIDs == "-1":
            self.medalsNeedShare = self.medalsLower20
            self.log.log("WARNING", "将分享所有等级小于20的直播间")
        else:
            try:
                self.medalsNeedShare = [
                    medal for medal in self.medalsLower20 if medal['medal']['target_id'] in
                    list(map(lambda x: int(x if x else 0), self.needShareUIDs.split(',')))
                ]
            except ValueError:
                self.medalsNeedShare = []
                self.log.log("ERROR", "需要分享的UID错误")

    async def likeInteract(self):
        '''
        点赞 *3 异步执行
        '''
        self.log.log("INFO", "点赞任务开始....(预计20秒完成)")
        likeTasks = [self.api.likeInteract(medal['room_info']['room_id']) for medal in self.medalsLower20]
        await asyncio.gather(*likeTasks)
        await asyncio.sleep(10)
        await self.getMedals()  # 刷新勋章
        self.log.log("SUCCESS", "点赞任务完成")
        finallyMedals = len(
            [medla for medla in self.medalsLower20 if medla['medal']['today_feed'] >= 600])
        msg = "20级以下牌子共 {} 个,完成点赞 {} 个".format(len(self.medalsLower20), finallyMedals)
        self.log.log("INFO", msg)
        if finallyMedals / len(self.medalsLower20) <= 0.8:
            self.log.log("WARNING", "点赞成功率过低,重新点赞任务")
            await self.likeInteract()

    async def shareRoom(self):
        '''
        分享直播间 CD 600s
        '''
        medalsNeedShare = self.medalsNeedShare.copy()
        if not medalsNeedShare:
            self.log.log("WARNING", "没有设置需要分享的直播间")
            return
        if len(medalsNeedShare) > 28:
            medalsNeedShare = medalsNeedShare[:28]
            self.log.log("WARNING", "由于B站分享CD为10分钟,所以一天最多只能分享28个直播间")
        needTime = len(medalsNeedShare) * 50 - 10
        self.log.log("INFO", "分享任务开始....(设置了 {} 个房间({})，预计{}分钟完成)".format(
            len(medalsNeedShare), "、".join([m['anchor_info']['nick_name'] for m in medalsNeedShare]), needTime))
        for index, medal in enumerate(medalsNeedShare):
            if medal['medal']['level'] >= 20:
                continue
            for i in range(1, 6):
                await self.api.shareRoom(medal['room_info']['room_id'])
                self.log.log("SUCCESS", "{} 分享成功 {} 次 (还需 {} 分钟完成)".format(
                    medal['anchor_info']['nick_name'], i, needTime))
                if i == 5:
                    break
                needTime -= 10
                await asyncio.sleep(600)
            if index == len(medalsNeedShare) - 1:
                break
            await asyncio.sleep(600)
        self.log.log("SUCCESS", "分享任务完成")

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
            finally:
                await asyncio.sleep(6)
        self.log.log("SUCCESS", "弹幕打卡任务完成")

    async def init(self):
        if not await self.loginVerify():
            self.log.log("ERROR", "登录失败")
            await self.session.close()
        else:
            await self.doSign()
            await self.getMedals()

    async def start(self):
        if self.isLogin:
            task = [self.likeInteract(), self.shareRoom(), self.sendDanmaku()]
            await asyncio.wait(task)
        await self.session.close()
